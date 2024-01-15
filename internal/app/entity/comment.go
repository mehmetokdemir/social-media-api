package entity

import "gorm.io/gorm"

// Comment DB Model
type Comment struct {
	gorm.Model
	UserID  uint   `gorm:"column:user_id"`
	User    User   `gorm:"foreignkey:UserID"`
	PostId  uint   `gorm:"column:post_id"`
	Post    Post   `gorm:"foreignkey:PostId"`
	Body    string `gorm:"column:body"`
	Image   string `gorm:"column:image"`
	ParenId *uint  `gorm:"column:parent_id"`
}
