package guard

import (
	"github.com/mehmetokdemir/social-media-api/internal/app/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IGuardRepository interface {
	CheckTokenInBlacklist(token string) bool
}

type guardRepository struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewRepository(db *gorm.DB, logger *zap.SugaredLogger) IGuardRepository {
	return &guardRepository{
		db:     db,
		logger: logger,
	}
}

func (r *guardRepository) CheckTokenInBlacklist(token string) bool {
	var count int64
	r.db.Model(&entity.BlackList{}).Where("token = ?", token).Count(&count)
	return count > 0
}
