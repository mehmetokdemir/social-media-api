package entity

import "gorm.io/gorm"

type ContentType string

const (
	ContentTypePost    ContentType = "post"
	ContentTypeComment ContentType = "comment"
)

type Like struct {
	gorm.Model
	UserID      uint        `gorm:"column:user_id"`
	User        User        `gorm:"foreignkey:UserID"`
	ContentType ContentType `gorm:"column:content_type"`

	ContentID uint `gorm:"column:content_id"` // Post or Comment ID
}
