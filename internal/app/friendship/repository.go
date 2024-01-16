package friendship

import (
	"errors"
	"github.com/mehmetokdemir/social-media-api/internal/app/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IFriendshipRepository interface {
	CreateFriendRequest(friendship entity.Friendship) (*entity.Friendship, error)
	AcceptFriendRequest(requestID uint) error
	DeleteFriendRequest(requestID uint) error
	ListFriendRequests(userID uint, status *entity.FriendshipStatusEnum) ([]*entity.Friendship, error)
	IsFriendShip(senderID, receiverID uint) (bool, error)
	IsFriendShipPending(senderID, receiverID uint) (bool, error)
	IsUserExist(userID uint) bool
	GetUserByID(userID uint) (*entity.User, error)
	GetFriendRequest(requestID uint) (*entity.Friendship, error)
	Migration() error
}

type friendshipRepository struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewRepository(db *gorm.DB, logger *zap.SugaredLogger) IFriendshipRepository {
	return &friendshipRepository{
		db:     db,
		logger: logger,
	}
}

func (r *friendshipRepository) CreateFriendRequest(friendship entity.Friendship) (*entity.Friendship, error) {
	if err := r.db.Create(&friendship).Error; err != nil {
		return nil, err
	}
	return &friendship, nil
}

func (r *friendshipRepository) AcceptFriendRequest(requestID uint) error {
	return r.db.Model(&entity.Friendship{}).Where("id =?", requestID).Update("status", entity.FriendshipStatusAccepted).Error
}

func (r *friendshipRepository) GetUserByID(id uint) (user *entity.User, err error) {
	if err = r.db.Model(&entity.User{}).Where("id =?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *friendshipRepository) IsFriendShip(senderID, receiverID uint) (bool, error) {
	var friendship entity.Friendship
	err := r.db.Model(&entity.Friendship{}).
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			senderID, receiverID, receiverID, senderID).
		Where("status = ?", entity.FriendshipStatusAccepted).
		First(&friendship).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}

func (r *friendshipRepository) IsFriendShipPending(senderID, receiverID uint) (bool, error) {
	var friendship entity.Friendship
	err := r.db.Model(&entity.Friendship{}).
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			senderID, receiverID, receiverID, senderID).
		Where("status = ?", entity.FriendshipStatusPending).
		First(&friendship).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}

func (r *friendshipRepository) GetFriends(userID uint, status *entity.FriendshipStatusEnum) ([]entity.Friendship, error) {
	var Friendships []entity.Friendship
	query := r.db.Where("(sender_id = ? OR receiver_id = ?)", userID, userID)
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Model(&Friendships).Preload("User").Find(&Friendships).Error; err != nil {
		return nil, err
	}

	return Friendships, nil
}

func (r *friendshipRepository) GetFriendshipRequest(senderID, receiverID uint) (Friendship *entity.Friendship, err error) {
	if err = r.db.Model(&entity.Friendship{}).Preload("User").
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			senderID, receiverID, receiverID, senderID).First(&Friendship).Error; err != nil {
		return nil, err
	}

	return Friendship, nil
}

func (r *friendshipRepository) IsUserExist(userID uint) bool {
	var user *entity.User
	if err := r.db.Model(&entity.User{}).Where("id =?", userID).First(&user).Error; err != nil {
		return false
	}
	return true
}

func (r *friendshipRepository) DeleteFriendRequest(requestID uint) error {
	return r.db.Delete(&entity.Friendship{}, requestID).Error
}

func (r *friendshipRepository) GetFriendRequest(requestID uint) (friendship *entity.Friendship, err error) {
	if err = r.db.Model(&entity.Friendship{}).Where("id =?", requestID).First(&friendship).Error; err != nil {
		return nil, err
	}
	return friendship, nil
}

func (r *friendshipRepository) ListFriendRequests(userID uint, status *entity.FriendshipStatusEnum) ([]*entity.Friendship, error) {
	friendships := make([]*entity.Friendship, 0)
	if status == nil {
		err := r.db.Model(&entity.Friendship{}).
			Where("sender_id = ? OR receiver_id = ?", userID, userID).
			Order("created_at DESC").
			Select("id, sender_id, receiver_id, status").
			Scan(&friendships).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := r.db.Model(&entity.Friendship{}).
			Where("sender_id = ? OR receiver_id = ?", userID, userID).
			Where("status = ?", *status).
			Order("created_at DESC").
			Select("id, sender_id, receiver_id, status").
			Scan(&friendships).Error
		if err != nil {
			return nil, err
		}
	}

	return friendships, nil
}

func (r *friendshipRepository) Migration() error {
	return r.db.AutoMigrate(entity.Friendship{})
}
