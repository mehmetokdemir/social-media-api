package friendship

import (
	"errors"
	"fmt"
	"github.com/mehmetokdemir/social-media-api/internal/app/common/httpmodel"
	"github.com/mehmetokdemir/social-media-api/internal/app/entity"
	"github.com/mehmetokdemir/social-media-api/internal/config"
	"go.uber.org/zap"
	"time"
)

type IFriendshipService interface {
	AddFriend(senderID, receiverID uint) error
	AcceptFriend(userID, friendshipRequestID uint) error
	RemoveFriend(userID, friendshipRequestID uint) error
	RejectFriend(userID, friendshipRequestID uint) error
	ListFriends(userID uint, status *entity.FriendshipStatusEnum) ([]ReadFriendship, error)
	GetFriendShipByRequestID(requestID uint) (*entity.Friendship, error)
}

type friendshipService struct {
	config               config.Config
	logger               *zap.SugaredLogger
	friendshipRepository IFriendshipRepository
}

func NewFriendshipService(friendshipRepository IFriendshipRepository, logger *zap.SugaredLogger, config config.Config) IFriendshipService {
	if friendshipRepository == nil {
		return nil
	}

	return &friendshipService{
		config:               config,
		friendshipRepository: friendshipRepository,
		logger:               logger,
	}
}

func (s *friendshipService) GetFriendShipByRequestID(requestID uint) (*entity.Friendship, error) {
	return s.friendshipRepository.GetFriendRequest(requestID)
}

func (s *friendshipService) AddFriend(senderID, receiverID uint) error {

	if senderID == receiverID {
		return errors.New("can not send yourself a friend request")
	}

	if !s.friendshipRepository.IsUserExist(senderID) || !s.friendshipRepository.IsUserExist(receiverID) {
		return errors.New("can not find users")
	}

	isFriend, err := s.friendshipRepository.IsFriendShip(senderID, receiverID)
	if err != nil {
		return err
	}

	isFriendShipPending, err := s.friendshipRepository.IsFriendShipPending(senderID, receiverID)
	if err != nil {
		return err
	}

	if isFriendShipPending || isFriend {
		return fmt.Errorf("already friend or waiting friend request from user %d", receiverID)
	}

	if _, err = s.friendshipRepository.CreateFriendRequest(entity.Friendship{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Status:     entity.FriendshipStatusPending,
	}); err != nil {
		return err
	}

	return nil
}

func (s *friendshipService) RemoveFriend(userID, friendshipRequestID uint) error {
	friendShip, err := s.friendshipRepository.GetFriendRequest(friendshipRequestID)
	if err != nil {
		return err
	}

	if friendShip.Status != entity.FriendshipStatusAccepted {
		return errors.New("can not delete friendship, because of you are not friends")
	}

	if friendShip.SenderID != userID || friendShip.ReceiverID != userID {
		return errors.New("do not have permission to remove this friend")
	}

	return s.friendshipRepository.DeleteFriendRequest(friendShip.ID)
}

func (s *friendshipService) RejectFriend(userID, friendshipRequestID uint) error {
	friendShip, err := s.friendshipRepository.GetFriendRequest(friendshipRequestID)
	if err != nil {
		return err
	}

	if friendShip.Status != entity.FriendshipStatusPending {
		return errors.New("can not reject friendship, if u want to reject first accept to friendship")
	}

	if friendShip.ReceiverID != userID {
		return errors.New("do not have permission to reject this friend")
	}

	return s.friendshipRepository.RejectFriendRequest(friendShip.ID)
}

func (s *friendshipService) ListFriends(userID uint, status *entity.FriendshipStatusEnum) ([]ReadFriendship, error) {
	friendShips, err := s.friendshipRepository.ListFriendRequests(userID, status)
	if err != nil {
		return nil, err
	}

	var readFriendships []ReadFriendship
	for _, fs := range friendShips {
		readFriendships = append(readFriendships, ReadFriendship{
			Id:        fs.ID,
			CreatedAt: fs.Model.CreatedAt.Format(time.RFC3339),
			Sender:    httpmodel.CommonUser{Id: fs.SenderID, Username: fs.Sender.Username, FirstName: fs.Sender.FirstName, LastName: fs.Sender.LastName, ProfilePhoto: fs.Sender.ProfilePhoto},
			Receiver:  httpmodel.CommonUser{Id: fs.ReceiverID, Username: fs.Receiver.Username, FirstName: fs.Receiver.FirstName, LastName: fs.Receiver.LastName, ProfilePhoto: fs.Receiver.ProfilePhoto},
			Status:    fs.Status,
		})
	}

	return readFriendships, nil
}

func (s *friendshipService) AcceptFriend(userID uint, friendshipRequestID uint) error {
	friendship, err := s.friendshipRepository.GetFriendRequest(friendshipRequestID)
	if err != nil {
		return err
	}

	if friendship.Status != entity.FriendshipStatusPending {
		return errors.New("can not accept to friendship, friendship status is not pending")
	}

	if friendship.ReceiverID != userID {
		return errors.New("can not accept to friendship")
	}

	return s.friendshipRepository.AcceptFriendRequest(friendshipRequestID)
}
