package post

import (
	"github.com/mehmetokdemir/social-media-api/internal/app/entity"
	"github.com/mehmetokdemir/social-media-api/internal/config"
	"go.uber.org/zap"
)

type IPostService interface {
	CreatePost(post entity.Post) (*entity.Post, error)
	UpdatePost(post entity.Post) (*entity.Post, error)
	GetPostById(id uint) (*entity.Post, error)
	DeletePostById(id uint) (*entity.Post, error)
	ListPosts() ([]entity.Post, error)
}

type postService struct {
	config     config.Config
	logger     *zap.SugaredLogger
	repository IPostRepository
}

func NewPostService(repository IPostRepository, logger *zap.SugaredLogger, config config.Config) IPostService {
	if repository == nil {
		return nil
	}

	return &postService{
		config:     config,
		repository: repository,
		logger:     logger,
	}
}

func (s *postService) CreatePost(post entity.Post) (*entity.Post, error) {
	return nil, nil
}

func (s *postService) UpdatePost(post entity.Post) (*entity.Post, error) {
	return nil, nil
}

func (s *postService) GetPostById(id uint) (*entity.Post, error) {
	return nil, nil
}

func (s *postService) DeletePostById(id uint) (*entity.Post, error) {
	return nil, nil
}

func (s *postService) ListPosts() ([]entity.Post, error) {
	return nil, nil
}
