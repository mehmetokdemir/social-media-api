package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	j "github.com/mehmetokdemir/social-media-api/internal/app/common/jwttoken"
	"github.com/mehmetokdemir/social-media-api/internal/app/entity"
	"github.com/mehmetokdemir/social-media-api/internal/app/user"
	"github.com/mehmetokdemir/social-media-api/internal/config"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type IAuthService interface {
	CreateToken(username, password string) (*LoginResponse, error)
	DeleteToken(token string) error
}

type authService struct {
	config              config.Config
	logger              *zap.SugaredLogger
	userService         user.IUserService
	blackListRepository IBlackListRepository
}

func NewAuthService(userService user.IUserService, blackListRepository IBlackListRepository, logger *zap.SugaredLogger, config config.Config) IAuthService {
	if blackListRepository == nil {
		return nil
	}

	return &authService{
		config:              config,
		userService:         userService,
		blackListRepository: blackListRepository,
		logger:              logger,
	}
}

func (s *authService) CreateToken(username, password string) (*LoginResponse, error) {

	userByUsername, err := s.userService.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if ok := s.verifyPassword(userByUsername.Password, password); !ok {
		return nil, errors.New("password mismatch")
	}

	tk := &j.Token{
		Username: username,
		UserId:   userByUsername.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(s.config.JwtATExpirationMinutes) * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	tokenString, err := token.SignedString([]byte(s.config.JwtATPrivateKey))
	if err != nil {
		return nil, errors.New("can not sign jwttoken")
	}

	if s.blackListRepository.CheckTokenInBlacklist(tokenString) {
		return nil, errors.New("invalid token")
	}

	return &LoginResponse{
		Username:  userByUsername.Username,
		Email:     userByUsername.Email,
		TokenHash: tokenString,
	}, nil
}

func (s *authService) verifyPassword(hashedPassword, requestedPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(requestedPassword)); err == nil {
		return true
	}
	return false
}

func (s *authService) DeleteToken(token string) error {
	return s.blackListRepository.CreateTokenToBlackList(entity.BlackList{
		Token: token,
	})
}
