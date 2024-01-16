package entity

import "gorm.io/gorm"

type FriendshipStatusEnum string

const (
	FriendshipStatusPending  FriendshipStatusEnum = "pending"
	FriendshipStatusAccepted FriendshipStatusEnum = "accepted"
)

// Friendship DB Model
type Friendship struct {
	gorm.Model
	SenderID   uint
	Sender     User `gorm:"foreignKey:SenderID"`
	ReceiverID uint
	Receiver   User `gorm:"foreignKey:ReceiverID"`
	Status     FriendshipStatusEnum
}
