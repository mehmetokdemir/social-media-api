package user

import (
	"github.com/mehmetokdemir/social-media-api/internal/app/entity"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	CreateUser(user entity.User) (*entity.User, error)
	CreateToken(username, password string) (string, error)
	VerifyPassword(hashedPassword, requestedPassword string) bool
	HashPassword(password string) (string, error)
}

type userService struct {
	logger         *zap.SugaredLogger
	userRepository IUserRepository
}

func NewUserService(userRepository IUserRepository, logger *zap.SugaredLogger) IUserService {
	return &userService{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (u *userService) CreateUser(user entity.User) (*entity.User, error) {
	return nil, nil
}

func (u *userService) CreateToken(username, password string) (string, error) {
	return "", nil
}

// TODO: Not need just use on the service for helper methods
func (u *userService) VerifyPassword(hashedPassword, requestedPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(requestedPassword)); err == nil {
		return true
	}
	return false
}

func (u *userService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
