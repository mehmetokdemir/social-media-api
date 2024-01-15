package entity

import "gorm.io/gorm"

// Post DB Model
type Post struct {
	gorm.Model
	UserID uint   `gorm:"column:user_id"`
	User   User   `gorm:"foreignkey:UserID"`
	Body   string `gorm:"column:body"`
	Image  string `gorm:"column:image"`
}
