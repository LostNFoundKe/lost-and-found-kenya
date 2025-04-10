package models

import (
	"github.com/google/uuid"
	"time"
)

// ItemStatus represents the status of a lost or found item
type ItemStatus string

const (
	ItemStatusLost     ItemStatus = "lost"
	ItemStatusFound    ItemStatus = "found"
	ItemStatusClaimed  ItemStatus = "claimed"
	ItemStatusReturned ItemStatus = "returned"
)

// Item represents a lost or found item
type Item struct {
	Model
	Title       string `gorm:"not null"`
	Description string `gorm:"type:text"`
	Category    string
	Status      ItemStatus `gorm:"not null;default:'lost'"`
	Location    string
	Date        time.Time
	Images      []Image
	UserID      uuid.UUID
	User        User
	Contact     string
	IsResolved  bool `gorm:"default:false"`
	Reward      float64
	Tags        []Tag `gorm:"many2many:item_tags;"`
}

// Tag represents a keyword associated with an item
type Tag struct {
	Model
	Name  string `gorm:"uniqueIndex;not null"`
	Items []Item `gorm:"many2many:item_tags;"`
}

// Claim represents a claim on a found item
type Claim struct {
	Model
	ItemID      uint
	ClaimerID   uuid.UUID
	Description string `gorm:"type:text"`
	Status      string `gorm:"default:'pending'"`
	ProofImages []ClaimImage
}

// ClaimImage represents proof images for a claim
type ClaimImage struct {
	Model
	URL     string `gorm:"not null"`
	ClaimID uint
}
