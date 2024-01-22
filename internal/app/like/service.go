package like

import (
	"fmt"
	"github.com/mehmetokdemir/social-media-api/internal/app/entity"
	"github.com/mehmetokdemir/social-media-api/internal/config"
	"go.uber.org/zap"
)

type ILikeService interface {
	DeleteLikesByCommentID(commentID uint) error
	DeleteLikesByPostID(postID uint) error
	GetCommentsLikeByID(commentID uint) ([]*entity.Like, error)
	GetPostsLikeByID(postID uint) ([]*entity.Like, error)
	GetPostCountLikeByID(postID uint) int64
	GetCommentCountLikeByID(commentID uint) int64
	LikePost(userID, postID uint) error
	LikeComment(userID, commentID uint) error
}

type likeService struct {
	config         config.Config
	logger         *zap.SugaredLogger
	likeRepository ILikeRepository
}

func NewLikeService(likeRepository ILikeRepository, logger *zap.SugaredLogger, config config.Config) ILikeService {
	if likeRepository == nil {
		return nil
	}
	return &likeService{
		config:         config,
		likeRepository: likeRepository,
		logger:         logger,
	}
}

func (s *likeService) DeleteLikesByCommentID(commentID uint) error {
	return s.likeRepository.DeleteLikes(commentID, entity.ContentTypeComment)
}

func (s *likeService) DeleteLikesByPostID(postID uint) error {
	return s.likeRepository.DeleteLikes(postID, entity.ContentTypePost)
}

func (s *likeService) LikePost(userID, postID uint) error {
	if ok := s.likeRepository.IsPostExist(postID); !ok {
		return fmt.Errorf("can not find post with id %d", postID)
	}

	isPostLikedByUser, err := s.likeRepository.IsPostLikedByUser(userID, postID)
	if err != nil {
		fmt.Println("err", err.Error())
		return err
	}

	if isPostLikedByUser {
		return fmt.Errorf("user already likes post which is id %d", postID)
	}

	if _, err = s.likeRepository.Like(entity.Like{
		UserID:      userID,
		ContentType: entity.ContentTypePost,
		ContentID:   postID,
	}); err != nil {
		return fmt.Errorf("can not like post because of %v", err.Error())
	}

	return nil
}

func (s *likeService) LikeComment(userID, commentID uint) error {
	if ok := s.likeRepository.IsCommentExist(commentID); !ok {
		return fmt.Errorf("can not find post with id %d", commentID)
	}

	if ok, err := s.likeRepository.IsCommentLikedByUser(userID, commentID); err != nil || ok {
		return fmt.Errorf("user already likes comment which is id %d", commentID)
	}

	if _, err := s.likeRepository.Like(entity.Like{
		UserID:      userID,
		ContentType: entity.ContentTypeComment,
		ContentID:   commentID,
	}); err != nil {
		return fmt.Errorf("can not like comment because of %v", err.Error())
	}

	return nil
}

func (s *likeService) GetCommentsLikeByID(commentID uint) ([]*entity.Like, error) {
	fmt.Println("get GetCommentsLikeByID")
	likes, err := s.likeRepository.GetLikesByID(commentID, entity.ContentTypeComment)
	if err != nil {
		fmt.Println("get GetCommentsLikeByID", err.Error())
		return nil, err
	}
	return likes, nil
}

func (s *likeService) GetPostCountLikeByID(postID uint) int64 {
	return s.likeRepository.GetCountOfLikes(postID, entity.ContentTypePost)
}
func (s *likeService) GetCommentCountLikeByID(commentID uint) int64 {
	return s.likeRepository.GetCountOfLikes(commentID, entity.ContentTypeComment)
}

func (s *likeService) GetPostsLikeByID(postID uint) ([]*entity.Like, error) {
	return s.likeRepository.GetLikesByID(postID, entity.ContentTypePost)
}
