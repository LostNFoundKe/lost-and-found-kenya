package models

type User struct {
	Model
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	FirstName string
	LastName  string
	Phone     string
	City      string
	IsAdmin   bool `gorm:"default:false"`
	Items     []Item
}
