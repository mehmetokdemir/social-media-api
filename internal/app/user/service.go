package user

import (
	"errors"
	"github.com/mehmetokdemir/social-media-api/internal/app/cdn"
	"github.com/mehmetokdemir/social-media-api/internal/app/entity"
	"github.com/mehmetokdemir/social-media-api/internal/config"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"mime/multipart"
)

type IUserService interface {
	CreateUser(user entity.User) (*entity.User, error)
	UpdateProfilePhoto(userID uint, file *multipart.FileHeader) (string, error)
	GetUserByUsername(username string) (*entity.User, error)
}

type userService struct {
	config         config.Config
	logger         *zap.SugaredLogger
	userRepository IUserRepository
	cdnService     cdn.ICdnService
}

func NewUserService(userRepository IUserRepository, cdnService cdn.ICdnService, logger *zap.SugaredLogger, config config.Config) IUserService {
	if userRepository == nil {
		return nil
	}

	return &userService{
		config:         config,
		userRepository: userRepository,
		logger:         logger,
		cdnService:     cdnService,
	}
}

func (s *userService) CreateUser(user entity.User) (*entity.User, error) {
	if s.userRepository.IsUserExistWithSameEmail(user.Email) || s.userRepository.IsUserExistWithSameUsername(user.Username) {
		return nil, errors.New("duplicated user")
	}

	hashedPassword, err := s.hashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword
	createdUser, err := s.userRepository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *userService) UpdateProfilePhoto(userID uint, file *multipart.FileHeader) (string, error) {
	userByID, err := s.userRepository.GetUserById(userID)
	if err != nil {
		return "", err
	}

	fileName, err := s.cdnService.UploadImage(file)
	if err != nil {
		return "", err
	}

	userByID.ProfilePhoto = fileName
	if _, err = s.userRepository.Update(*userByID); err != nil {
		return "", err
	}

	return fileName, err
}

func (s *userService) GetUserByUsername(username string) (*entity.User, error) {
	return s.userRepository.GetUserByUsername(username)
}

func (s *userService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
