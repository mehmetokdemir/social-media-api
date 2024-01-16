package auth

import (
	"github.com/mehmetokdemir/social-media-api/internal/app/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IBlackListRepository interface {
	CreateTokenToBlackList(blackList entity.BlackList) error
	CheckTokenInBlacklist(token string) bool
	Migration() error
}

type blackListRepository struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewRepository(db *gorm.DB, logger *zap.SugaredLogger) IBlackListRepository {
	return &blackListRepository{
		db:     db,
		logger: logger,
	}
}

func (r *blackListRepository) FindAll() ([]*entity.BlackList, error) {
	var blacklist []*entity.BlackList
	err := r.db.Find(&blacklist).Error
	return blacklist, err
}

func (r *blackListRepository) CreateTokenToBlackList(blackList entity.BlackList) error {
	if err := r.db.Create(&blackList).Error; err != nil {
		return err
	}
	return nil
}

func (r *blackListRepository) CheckTokenInBlacklist(token string) bool {
	var count int64
	r.db.Model(&entity.BlackList{}).Where("token = ?", token).Count(&count)
	return count > 0
}

func (r *blackListRepository) Migration() error {
	return r.db.AutoMigrate(entity.BlackList{})
}
