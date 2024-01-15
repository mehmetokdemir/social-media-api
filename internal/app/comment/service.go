package comment

import (
	"github.com/mehmetokdemir/social-media-api/internal/app/entity"
	"github.com/mehmetokdemir/social-media-api/internal/config"
	"go.uber.org/zap"
)

type ICommentService interface {
	CreateComment(comment entity.Comment) (*entity.Comment, error)
	UpdateComment(comment entity.Comment) (*entity.Comment, error)
	GetCommentById(id uint) (*entity.Comment, error)
	DeleteCommentById(id uint) error
	ListComments(postID uint) ([]entity.Comment, error)
}

type commentService struct {
	config     config.Config
	logger     *zap.SugaredLogger
	repository ICommentRepository
}

func NewCommentService(repository ICommentRepository, logger *zap.SugaredLogger, config config.Config) ICommentService {
	if repository == nil {
		return nil
	}

	return &commentService{
		config:     config,
		repository: repository,
		logger:     logger,
	}
}

func (s *commentService) CreateComment(comment entity.Comment) (*entity.Comment, error) {
	return nil, nil
}

func (s *commentService) UpdateComment(comment entity.Comment) (*entity.Comment, error) {
	return nil, nil
}

func (s *commentService) GetCommentById(id uint) (*entity.Comment, error) {
	return nil, nil
}

func (s *commentService) DeleteCommentById(id uint) error {
	return nil
}

func (s *commentService) ListComments(postID uint) ([]entity.Comment, error) {
	return nil, nil
}
