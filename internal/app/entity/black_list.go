package entity

type BlackList struct {
	Token string `gorm:"primaryKey;autoIncrement:false"`
}
