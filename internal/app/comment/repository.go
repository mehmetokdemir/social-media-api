package comment

import (
	"errors"
	"github.com/mehmetokdemir/social-media-api/internal/app/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ICommentRepository interface {
	Create(comment entity.Comment) (*entity.Comment, error)
	Update(comment entity.Comment) (*entity.Comment, error)
	Get(id uint) (*entity.Comment, error)
	Delete(id uint) error
	DeleteCommentsByPostID(postID uint) error
	List() ([]*entity.Comment, error)
	ListCommentsByPostID(postID uint) ([]entity.Comment, error)
	ListMainCommentsByPostID(postID uint) ([]entity.Comment, error)
	IsPostExist(postID uint) bool

	ListCommentsByParentID(parentCommentID uint) ([]entity.Comment, error)
	DeleteCommentsByParentID(parentCommentID uint) error

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

func (r *commentRepository) Create(comment entity.Comment) (*entity.Comment, error) {
	if err := r.db.Create(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) Update(comment entity.Comment) (*entity.Comment, error) {
	if err := r.db.Model(entity.Comment{}).Where("id =?", comment.ID).Updates(comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) Delete(id uint) error {
	err := r.db.Model(&entity.Comment{}).Where("id = ? ", id).Delete(nil).Error
	if err != nil {
		return err
	}

	if r.db.RowsAffected != 0 {
		return errors.New("can not find comments to be deleted")
	}

	return nil
}

func (r *commentRepository) DeleteCommentsByParentID(id uint) error {
	result := r.db.Where("parent_id = ? ", id).Delete(&entity.Comment{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}

func (r *commentRepository) List() ([]*entity.Comment, error) {
	var comments []*entity.Comment
	if err := r.db.Model(&entity.Comment{}).Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *commentRepository) ListCommentsByParentID(commentParentID uint) ([]entity.Comment, error) {
	var comments []entity.Comment
	if err := r.db.Preload("User").Model(&entity.Comment{}).Where("parent_id = ? ", commentParentID).Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *commentRepository) ListCommentsByPostID(postID uint) ([]entity.Comment, error) {
	var comments []entity.Comment
	if err := r.db.Preload("User").Model(&entity.Comment{}).Where("post_id = ?", postID).Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *commentRepository) ListMainCommentsByPostID(postID uint) ([]entity.Comment, error) {
	var comments []entity.Comment
	if err := r.db.Preload("User").Model(&entity.Comment{}).Where("post_id = ? AND parent_id IS NULL", postID).Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *commentRepository) IsPostExist(postID uint) bool {
	var post *entity.Post
	if err := r.db.Model(&entity.Post{}).Where("id =?", postID).First(&post).Error; err != nil || post == nil {
		return false
	}
	return true
}

func (r *commentRepository) Get(id uint) (comment *entity.Comment, err error) {
	if err = r.db.Model(&entity.Comment{}).Where("id =?", id).First(&comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

func (r *commentRepository) DeleteCommentsByPostID(postID uint) error {
	result := r.db.Where("post_id = ?", postID).Delete(&entity.Comment{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}

func (r *commentRepository) Migration() error {
	return r.db.AutoMigrate(entity.Comment{})
}
