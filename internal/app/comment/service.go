package comment

import (
	"errors"
	"fmt"
	"github.com/mehmetokdemir/social-media-api/internal/app/cdn"
	"github.com/mehmetokdemir/social-media-api/internal/app/entity"
	"github.com/mehmetokdemir/social-media-api/internal/app/like"
	"github.com/mehmetokdemir/social-media-api/internal/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mime/multipart"
	"time"
)

type ICommentService interface {
	CreateComment(userID uint, comment CreateRequest) (*entity.Comment, error)

	UpdateComment(userID uint, req UpdateRequest) (*entity.Comment, error)
	UpdateCommentImage(commentID, userID uint, header *multipart.FileHeader) (string, error)

	GetCommentById(id uint) (*entity.Comment, error)
	DeleteCommentById(userID, id uint) error

	ListCommentsByPostID(postID uint) ([]*entity.Comment, error)
	DeleteCommentsByPostID(postID uint) error
}

type commentService struct {
	config      config.Config
	logger      *zap.SugaredLogger
	repository  ICommentRepository
	cdnService  cdn.ICdnService
	likeService like.ILikeService
}

func NewCommentService(repository ICommentRepository, likeService like.ILikeService, cdnService cdn.ICdnService, logger *zap.SugaredLogger, config config.Config) ICommentService {
	if repository == nil {
		return nil
	}

	return &commentService{
		config:      config,
		repository:  repository,
		logger:      logger,
		cdnService:  cdnService,
		likeService: likeService,
	}
}

func (s *commentService) CreateComment(userID uint, comment CreateRequest) (*entity.Comment, error) {
	if ok := s.repository.IsPostExist(comment.PostId); !ok {
		return nil, errors.New("post not found")
	}

	dbComment := entity.Comment{
		UserID: userID,
		PostID: comment.PostId,
		Body:   comment.Body,
	}

	if comment.ParentID != nil {
		fmt.Println("girdi", *comment.ParentID)
		if _, err := s.repository.Get(*comment.ParentID); err != nil {
			return nil, errors.New("parent comment not found")
		}
		dbComment.ParenId = comment.ParentID
	}

	return s.repository.Create(dbComment)
}

func (s *commentService) UpdateComment(userID uint, comment UpdateRequest) (*entity.Comment, error) {
	commentByID, err := s.repository.Get(comment.Id)
	if err != nil {
		return nil, err
	}

	if commentByID.UserID != userID {
		return nil, errors.New("do not have permission to update this comment")
	}

	if time.Since(commentByID.Model.CreatedAt) > 5 {
		return nil, errors.New("comment update period has expired")
	}

	return s.repository.Update(entity.Comment{
		Model:   gorm.Model{ID: comment.Id},
		UserID:  userID,
		PostID:  commentByID.PostID,
		Body:    comment.Body,
		Image:   commentByID.Image,
		ParenId: commentByID.ParenId,
	})
}

func (s *commentService) UpdateCommentImage(commentID, userID uint, file *multipart.FileHeader) (string, error) {
	commentByID, err := s.repository.Get(commentID)
	if err != nil {
		return "", err
	}

	if commentByID.UserID != userID {
		return "", errors.New("do not have permission to update this post")
	}

	fileName, err := s.cdnService.UploadImage(file)
	if err != nil {
		return "", err
	}

	commentByID.Image = fileName
	if _, err = s.repository.Update(*commentByID); err != nil {
		return "", err
	}

	return fileName, err
}

func (s *commentService) GetCommentById(id uint) (*entity.Comment, error) {
	return s.repository.Get(id)
}

func (s *commentService) DeleteCommentById(userID, id uint) error {

	commentById, err := s.repository.Get(id)
	if err != nil {
		return err
	}

	if userID != commentById.UserID {
		return errors.New("do not have permission to delete this comment")
	}

	if err = s.likeService.DeleteLikesByCommentID(commentById.ID); err != nil {
		return err
	}

	if commentById.ParenId == nil {
		if err = s.repository.DeleteCommentsWithSubComments(commentById.ID); err != nil {
			return err
		}
	} else {
		if err = s.repository.Delete(commentById.ID); err != nil {
			return err
		}
	}

	return nil
}

func (s *commentService) DeleteCommentsByPostID(postID uint) error {
	comments, err := s.repository.ListCommentsByPostID(postID)
	if err != nil {
		return err
	}

	if len(comments) > 0 {
		return s.repository.DeleteCommentsByPostID(postID)
	}
	return nil
}

func (s *commentService) ListCommentsByPostID(postID uint) ([]*entity.Comment, error) {
	return s.repository.ListCommentsByPostID(postID)
}
