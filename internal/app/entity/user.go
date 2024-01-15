package entity

import "gorm.io/gorm"

// User DB Model
type User struct {
	gorm.Model
	FirstName    string `gorm:"column:first_name"`
	LastName     string `gorm:"column:last_name"`
	Email        string `gorm:"uniqueIndex"`
	Username     string `gorm:"uniqueIndex"`
	Password     string `gorm:"column:password"`
	PhoneNumber  string `gorm:"column:phone_number"`
	ProfilePhoto string `gorm:"column:profile_photo"`
}
