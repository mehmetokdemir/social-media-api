package friendship

import (
	"github.com/mehmetokdemir/social-media-api/internal/app/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IFriendshipRepository interface {
	CreateFriendRequest(Friendship entity.Friendship) (*entity.Friendship, error)
	UpdateFriendshipStatus(requestID uint, status entity.FriendshipStatusEnum) error
	DeleteFriendship(requestID uint) error
	GetFriendshipRequest(senderID, receiverID uint) (*entity.Friendship, error)
	GetFriends(userID uint, status *entity.FriendshipStatusEnum) ([]entity.Friendship, error)
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

func (r *friendshipRepository) CreateFriendRequest(Friendship entity.Friendship) (*entity.Friendship, error) {
	if err := r.db.Create(&Friendship).Error; err != nil {
		return nil, err
	}
	return &Friendship, nil
}

func (r *friendshipRepository) UpdateFriendshipStatus(requestID uint, status entity.FriendshipStatusEnum) error {
	return r.db.Model(&entity.Friendship{}).Where("ID =?", requestID).Update("status", status).Error
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
	if err = r.db.Model(&entity.Friendship{}).Preload("User").Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		senderID, receiverID, receiverID, senderID).First(&Friendship).Error; err != nil {
		return nil, err
	}

	return Friendship, nil
}

func (r *friendshipRepository) DeleteFriendship(requestID uint) error {
	return r.db.Delete(&entity.Friendship{}, requestID).Error
}

func (r *friendshipRepository) Migration() error {
	return r.db.AutoMigrate(entity.Friendship{})
}
