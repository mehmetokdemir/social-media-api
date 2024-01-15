package comment

import (
	"github.com/mehmetokdemir/social-media-api/internal/app/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ICommentRepository interface {
	Create(post entity.Comment) (*entity.Comment, error)
	Update(post entity.Comment) (*entity.Comment, error)
	Delete(id uint) error
	List() ([]entity.Comment, error)
	IsPostExist(postID uint) bool
	Migration() error
}

type commentRepository struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewRepository(db *gorm.DB, logger *zap.SugaredLogger) ICommentRepository {
	return &commentRepository{
		db:     db,
		logger: logger,
	}
}

func (r *commentRepository) Create(post entity.Comment) (*entity.Comment, error) {
	if err := r.db.Create(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *commentRepository) Update(post entity.Comment) (*entity.Comment, error) {
	if err := r.db.Create(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *commentRepository) Delete(id uint) error {
	return nil
}

func (r *commentRepository) List() ([]entity.Comment, error) {
	return nil, nil
}

func (r *commentRepository) IsPostExist(postID uint) bool {
	var post *entity.Post
	if err := r.db.Model(&entity.Post{}).Where("ID =?", postID).First(&post).Error; err != nil || post == nil {
		return false
	}

	return true
}

func (r *commentRepository) Migration() error {
	return r.db.AutoMigrate(entity.Comment{})
}
