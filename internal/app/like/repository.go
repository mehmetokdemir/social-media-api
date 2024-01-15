package like

import (
	"github.com/mehmetokdemir/social-media-api/internal/app/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ILikeRepository interface {
	// todo: burdan devam

	Create(post entity.Post) (*entity.Post, error)
	Migration() error
}

type likeRepository struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewRepository(db *gorm.DB, logger *zap.SugaredLogger) ILikeRepository {
	return &likeRepository{
		db:     db,
		logger: logger,
	}
}

func (r *likeRepository) Create(post entity.Post) (*entity.Post, error) {
	if err := r.db.Create(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *likeRepository) Migration() error {
	return r.db.AutoMigrate(entity.Like{})
}
