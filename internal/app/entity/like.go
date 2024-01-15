package entity

import "gorm.io/gorm"

type LikeType string

const (
	PostLike    LikeType = "post"
	CommentLike LikeType = "comment"
)

// TODO: not sure ?
type Like struct {
	gorm.Model
	UserID          uint `gorm:"column:user_id"`
	User            User `gorm:"foreignkey:UserID"`
	Type            LikeType
	PostOrCommentID uint
}
