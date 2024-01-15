package post

import (
	"github.com/mehmetokdemir/social-media-api/internal/app/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IPostRepository interface {
	Create(post entity.Post) (*entity.Post, error)
	Update(post entity.Post) (*entity.Post, error)
	Delete(id uint) error
	List() ([]entity.Post, error)
	Migration() error
}

type postRepository struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewRepository(db *gorm.DB, logger *zap.SugaredLogger) IPostRepository {
	return &postRepository{
		db:     db,
		logger: logger,
	}
}

func (r *postRepository) Create(post entity.Post) (*entity.Post, error) {
	if err := r.db.Create(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) Update(post entity.Post) (*entity.Post, error) {
	if err := r.db.Create(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) Delete(id uint) error {
	return nil
}

func (r *postRepository) List() ([]entity.Post, error) {
	return nil, nil
}

func (r *postRepository) Migration() error {
	return r.db.AutoMigrate(entity.Post{})
}
