package post

import (
	"fmt"
	"github.com/mehmetokdemir/social-media-api/internal/app/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IPostRepository interface {
	Create(post entity.Post) (*entity.Post, error)
	Update(post entity.Post) (*entity.Post, error)
	Get(id uint) (*entity.Post, error)
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
	if err := r.db.Model(entity.Post{}).Where("id =?", post.ID).Updates(post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) Delete(id uint) error {
	fmt.Println("delete DeleteRepoPost")
	return r.db.Where("id = ?", id).Delete(&entity.Post{}).Error
}

func (r *postRepository) Get(id uint) (post *entity.Post, err error) {
	fmt.Println("get get endpoint")
	if err = r.db.Model(&entity.Post{}).Where("id =?", id).First(&post).Error; err != nil {
		fmt.Println("get get endpoint err", err.Error())
		return nil, err
	}
	return post, nil
}

func (r *postRepository) List() ([]entity.Post, error) {
	var posts []entity.Post
	if err := r.db.Model(&entity.Post{}).Order("created_at DESC").Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *postRepository) Migration() error {
	return r.db.AutoMigrate(entity.Post{})
}
