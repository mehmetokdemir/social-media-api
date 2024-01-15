package friendship

import (
	"github.com/mehmetokdemir/social-media-api/internal/app/user"
	"github.com/mehmetokdemir/social-media-api/internal/config"
	"go.uber.org/zap"
)

type IFriendshipService interface {
	AddFriend(senderID, receiverID uint) error
	RemoveFriend(FriendshipRequestID uint) error
}

type FriendshipService struct {
	config               config.Config
	logger               *zap.SugaredLogger
	FriendshipRepository IFriendshipRepository
	userRepository       user.IUserRepository
}

func NewFriendshipService(userRepository user.IUserRepository, FriendshipRepository IFriendshipRepository, logger *zap.SugaredLogger, config config.Config) IFriendshipService {
	if userRepository == nil || FriendshipRepository == nil {
		return nil
	}

	return &FriendshipService{
		config:               config,
		userRepository:       userRepository,
		FriendshipRepository: FriendshipRepository,
		logger:               logger,
	}
}

func (s *FriendshipService) AddFriend(senderID, receiverID uint) error {
	return nil
}

func (s *FriendshipService) RemoveFriend(FriendshipRequestID uint) error {
	return nil
}
