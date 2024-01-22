package like

import (
	"github.com/mehmetokdemir/social-media-api/internal/app/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ILikeRepository interface {
	Like(like entity.Like) (*entity.Like, error)

	IsPostExist(postID uint) bool
	IsCommentExist(commentID uint) bool

	IsPostLikedByUser(userID uint, postID uint) (bool, error)
	IsCommentLikedByUser(userID uint, commentID uint) (bool, error)

	GetLikesByID(contentID uint, contentType entity.ContentType) ([]*entity.Like, error)
	DeleteLikes(id uint, contentType entity.ContentType) error
	GetCountOfLikes(id uint, contentType entity.ContentType) int64

	ListCommentLikesByParentID(parentID uint) ([]*entity.Like, error)
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

func (r *likeRepository) Like(like entity.Like) (*entity.Like, error) {
	if err := r.db.Create(&like).Error; err != nil {
		return nil, err
	}
	return &like, nil
}

func (r *likeRepository) IsPostExist(postID uint) bool {
	var post *entity.Post
	if err := r.db.Model(&entity.Post{}).Where("id =?", postID).First(&post).Error; err != nil || post == nil {
		return false
	}

	return true
}

func (r *likeRepository) IsCommentExist(commentID uint) bool {
	var comment *entity.Comment
	if err := r.db.Model(&entity.Comment{}).Where("id =?", commentID).First(&comment).Error; err != nil || comment == nil {
		return false
	}

	return true
}

func (r *likeRepository) IsPostLikedByUser(userID, postID uint) (bool, error) {
	var like entity.Like
	if err := r.db.Where("user_id = ? AND content_type = ? AND content_id = ?", userID, entity.ContentTypePost, postID).First(&like).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return false, nil
		default:
			return true, err
		}
	}

	return true, nil
}

func (r *likeRepository) IsCommentLikedByUser(userID, commentID uint) (bool, error) {
	var like entity.Like
	if err := r.db.Where("user_id = ? AND content_type = ? AND content_id = ?", userID, entity.ContentTypeComment, commentID).First(&like).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return false, nil
		default:
			return true, err
		}
	}

	return true, nil
}

func (r *likeRepository) GetLikesByID(contentID uint, contentType entity.ContentType) ([]*entity.Like, error) {
	var likes []*entity.Like
	if err := r.db.Model(&entity.Like{}).Where("content_id = ? AND content_type = ?", contentID, contentType).Find(&likes).Error; err != nil {
		return nil, err
	}
	return likes, nil
}

func (r *likeRepository) DeleteLikes(id uint, contentType entity.ContentType) error {
	result := r.db.Where("content_id = ? AND content_type = ?", id, contentType).Delete(&entity.Like{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}

func (r *likeRepository) ListCommentLikesByParentID(parentID uint) ([]*entity.Like, error) {
	var likes []*entity.Like
	if err := r.db.Model(&entity.Like{}).Where("parent_id = ?", parentID).Find(&likes).Error; err != nil {
		return nil, err
	}

	return likes, nil
}

func (r *likeRepository) GetCountOfLikes(id uint, contentType entity.ContentType) int64 {
	var count int64
	if err := r.db.Model(&entity.Like{}).Where("content_id = ? AND content_type = ?", id, contentType).Count(&count).Error; err != nil {
		return 0
	}
	return count
}

func (r *likeRepository) Migration() error {
	return r.db.AutoMigrate(entity.Like{})
}
