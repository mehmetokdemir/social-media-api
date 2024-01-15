package user

import (
	"fmt"
	"github.com/mehmetokdemir/social-media-api/internal/app/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(user entity.User) (*entity.User, error)
	UpdateProfilePhoto(photo string) error
	GetUserByUsername(username string) (*entity.User, error)
	GetUserById(id uint) (*entity.User, error)
	IsUserExistWithSameUsername(username string) bool
	IsUserExistWithSameEmail(email string) bool
	Migration() error
}

type userRepository struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewRepository(db *gorm.DB, logger *zap.SugaredLogger) IUserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

func (r *userRepository) CreateUser(user entity.User) (*entity.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserById(id uint) (user *entity.User, err error) {
	if err = r.db.Model(&entity.User{}).Where("ID =?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserByUsername(username string) (user *entity.User, err error) {
	if err = r.db.Model(&entity.User{}).Where("username =?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) IsUserExistWithSameEmail(email string) bool {
	return r.isUserExistWithCredential("email", email)
}

func (r *userRepository) isUserExistWithCredential(key, value string) bool {
	var user *entity.User
	if err := r.db.Where(fmt.Sprintf("%s =?", key), value).First(&user).Error; err == nil && user != nil {
		return true
	}
	return false
}

func (r *userRepository) IsUserExistWithSameUsername(username string) bool {
	return r.isUserExistWithCredential("username", username)
}

func (r *userRepository) UpdateProfilePhoto(photo string) error {
	return nil
}

func (r *userRepository) Migration() error {
	return r.db.AutoMigrate(entity.User{})
}
