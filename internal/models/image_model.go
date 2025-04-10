package models

import "github.com/google/uuid"

// Image represents an image of a lost or found item
type Image struct {
	Model
	URL    string `gorm:"not null"`
	ItemID uuid.UUID
}
