package post

import (
	"errors"
	"fmt"
	"github.com/mehmetokdemir/social-media-api/internal/app/cdn"
	"github.com/mehmetokdemir/social-media-api/internal/app/comment"
	"github.com/mehmetokdemir/social-media-api/internal/app/common/httpmodel"
	"github.com/mehmetokdemir/social-media-api/internal/app/entity"
	"github.com/mehmetokdemir/social-media-api/internal/app/like"
	"github.com/mehmetokdemir/social-media-api/internal/app/transaction"
	"github.com/mehmetokdemir/social-media-api/internal/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mime/multipart"
	"time"
)

type IPostService interface {
	CreatePost(post entity.Post) (*entity.Post, error)
	UpdatePost(userID uint, post UpdateRequest) (*entity.Post, error)
	GetPostById(id uint) (*ReadPostResponse, error)
	DeletePostById(userID uint, id uint) error
	ListPosts() ([]ReadPostResponse, error)
	UpdatePostImage(postID, userID uint, header *multipart.FileHeader) (string, error)
}

type postService struct {
	config             config.Config
	logger             *zap.SugaredLogger
	repository         IPostRepository
	commentService     comment.ICommentService
	likeService        like.ILikeService
	transactionService transaction.ITransactionService
	cdnService         cdn.ICdnService
}

func NewPostService(repository IPostRepository, cdnService cdn.ICdnService, transactionService transaction.ITransactionService, commentService comment.ICommentService, likeService like.ILikeService, logger *zap.SugaredLogger, config config.Config) IPostService {
	if repository == nil {
		return nil
	}

	return &postService{
		config:             config,
		repository:         repository,
		transactionService: transactionService,
		commentService:     commentService,
		likeService:        likeService,
		logger:             logger,
		cdnService:         cdnService,
	}
}

func (s *postService) CreatePost(post entity.Post) (*entity.Post, error) {
	return s.repository.Create(post)
}

func (s *postService) UpdatePostImage(postID, userID uint, file *multipart.FileHeader) (string, error) {
	postById, err := s.repository.Get(postID)
	if err != nil {
		return "", err
	}

	if postById.UserID != userID {
		return "", errors.New("do not have permission to update this post")
	}

	fileName, err := s.cdnService.UploadImage(file)
	if err != nil {
		return "", err
	}

	postById.Image = fileName
	if _, err = s.repository.Update(*postById); err != nil {
		return "", err
	}

	return fileName, err
}

func (s *postService) UpdatePost(userID uint, post UpdateRequest) (*entity.Post, error) {
	postById, err := s.repository.Get(post.Id)
	if err != nil {
		return nil, err
	}

	if postById.UserID != userID {
		return nil, errors.New("do not have permission to update this post")
	}

	if time.Since(postById.Model.CreatedAt).Minutes() > 5 {
		return nil, errors.New("post update period has expired")
	}

	post.Image = postById.Image

	return s.repository.Update(entity.Post{
		Model:  gorm.Model{ID: post.Id},
		UserID: userID,
		Body:   post.Body,
		Image:  post.Image,
	})
}

func (s *postService) GetPostById(id uint) (*ReadPostResponse, error) {
	post, err := s.repository.Get(id)
	if err != nil {
		return nil, err
	}

	return &ReadPostResponse{
		Id:        post.ID,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
		Body:      post.Body,
		Image:     post.Image,
	}, nil
}

func (s *postService) DeletePostById(userID uint, id uint) error {
	// TODO : Need transaction commit with rollback
	postByID, err := s.repository.Get(id)
	if err != nil {
		return err
	}

	if postByID.UserID != userID {
		return errors.New("do not have permission to update this post")
	}

	// Get comments which is belongs to post
	comments, err := s.commentService.ListCommentsByPostID(id)
	if err != nil {
		return err
	}

	for _, com := range comments {
		likes, err := s.likeService.GetCommentsLikeByID(com.ID)
		if err == nil && len(likes) > 0 {
			for _, l := range likes {
				if err = s.likeService.DeleteLikesByCommentID(l.ContentID); err != nil {
					fmt.Println("girdi 3", err.Error())
					return err
				}
			}
		}
	}

	if err = s.commentService.DeleteCommentsByPostID(id); err != nil {
		return err
	}

	if err = s.likeService.DeleteLikesByPostID(id); err != nil {
		return err
	}

	if err = s.repository.Delete(id); err != nil {
		return err
	}

	if err = s.transactionService.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *postService) ListPosts() ([]ReadPostResponse, error) {
	posts, err := s.repository.List()
	if err != nil {
		return nil, err
	}

	var rsp []ReadPostResponse
	for _, post := range posts {
		rsp = append(rsp, ReadPostResponse{
			Id:        post.ID,
			CreatedAt: post.CreatedAt.Format(time.RFC3339),
			Body:      post.Body,
			Image:     post.Image,
			User:      httpmodel.CommonUser{Id: post.UserID, Username: post.User.Username, ProfilePhoto: post.User.ProfilePhoto},
		})
	}

	return rsp, nil
}
