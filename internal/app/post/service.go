package post

import (
	"errors"
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

	comments, err := s.commentService.ListCommentsByPostID(postByID.ID)
	if err != nil {
		return err
	}

	for _, com := range comments {
		// Delete comments likes which are belongs to post
		if err = s.likeService.DeleteLikesByCommentID(com.ID); err != nil {
			return err
		}
	}

	// Delete post likes by post id
	if err = s.likeService.DeleteLikesByPostID(postByID.ID); err != nil {
		return err
	}

	// Delete comments by post id
	if err = s.commentService.DeleteCommentsByPostID(postByID.ID); err != nil {
		return err
	}

	// Delete post by id
	if err = s.repository.Delete(postByID.ID); err != nil {
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
		var rspComments []ReadPostResponseComment
		// Get main comments without belongs to other comments
		comments, err := s.commentService.ListMainCommentsByPostID(post.ID)
		if err != nil {
			return nil, err
		}

		for _, com := range comments {
			var rspSubComments []ReadPostResponseComment
			subComments, err := s.commentService.ListCommentsByParentID(com.ID)
			if err != nil {
				return nil, err
			}

			for _, subComment := range subComments {
				rspSubComments = append(rspSubComments, ReadPostResponseComment{
					Id:         subComment.ID,
					Body:       subComment.Body,
					Image:      subComment.Image,
					User:       httpmodel.CommonUser{Id: subComment.UserID, Username: subComment.User.Username, FirstName: subComment.User.FirstName, LastName: subComment.User.LastName, ProfilePhoto: subComment.User.ProfilePhoto},
					LikedCount: s.likeService.GetCommentCountLikeByID(subComment.ID),
				})
			}

			rspComments = append(rspComments, ReadPostResponseComment{
				Id:          com.ID,
				Body:        com.Body,
				Image:       com.Image,
				User:        httpmodel.CommonUser{Id: com.UserID, Username: com.User.Username, FirstName: com.User.FirstName, LastName: com.User.LastName, ProfilePhoto: com.User.ProfilePhoto},
				LikedCount:  s.likeService.GetCommentCountLikeByID(com.ID),
				SubComments: rspSubComments,
			})
		}

		rsp = append(rsp, ReadPostResponse{
			Id:         post.ID,
			CreatedAt:  post.CreatedAt.Format(time.RFC3339),
			User:       httpmodel.CommonUser{Id: post.UserID, Username: post.User.Username, FirstName: post.User.FirstName, LastName: post.User.LastName, ProfilePhoto: post.User.ProfilePhoto},
			Body:       post.Body,
			Image:      post.Image,
			LikedCount: s.likeService.GetPostCountLikeByID(post.ID),
			Comments:   rspComments,
		})
	}

	return rsp, nil
}
