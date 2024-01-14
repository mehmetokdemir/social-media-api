package user

import (
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

func (u *userRepository) CreateUser(user entity.User) (*entity.User, error) {
	return nil, nil
}

func (u *userRepository) GetUserById(id uint) (*entity.User, error) {
	return nil, nil
}

func (u *userRepository) GetUserByUsername(username string) (*entity.User, error) {
	return nil, nil
}

func (u *userRepository) IsUserExistWithSameUsername(username string) bool {
	return false
}

func (u *userRepository) IsUserExistWithSameEmail(email string) bool {
	return false
}

func (u *userRepository) UpdateProfilePhoto(photo string) error {
	return nil
}

func (u *userRepository) Migration() error {
	return u.db.AutoMigrate(entity.User{})
}
