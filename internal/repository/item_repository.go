package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"lostnfound-api/internal/models"
)

// ItemRepository handles database operations for items
type ItemRepository struct {
	db *gorm.DB
}

// NewItemRepository creates a new ItemRepository
func NewItemRepository(db *gorm.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

// Create adds a new item to the database
func (r *ItemRepository) Create(item *models.Item) error {
	return r.db.Create(item).Error
}

// GetByID retrieves an item by ID
func (r *ItemRepository) GetByID(id uuid.UUID) (*models.Item, error) {
	var item models.Item
	err := r.db.Preload("Images").Preload("User").Preload("Tags").First(&item, id).Error
	return &item, err
}

// List retrieves items with filtering options
func (r *ItemRepository) List(status string, category string, page, limit int) ([]models.Item, int64, error) {
	var items []models.Item
	var count int64

	query := r.db.Model(&models.Item{})

	// Apply filters
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// Get total count
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	err = query.Preload("Images").Preload("User").Offset(offset).Limit(limit).Order("created_at DESC").Find(&items).Error

	return items, count, err
}

// Update updates an existing item
func (r *ItemRepository) Update(item *models.Item) error {
	return r.db.Save(item).Error
}

// Delete removes an item from the database
func (r *ItemRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Item{}, id).Error
}

// SearchByKeyword searches items by keyword
func (r *ItemRepository) SearchByKeyword(keyword string, page, limit int) ([]models.Item, int64, error) {
	var items []models.Item
	var count int64

	query := r.db.Model(&models.Item{}).
		Where("title ILIKE ?", "%"+keyword+"%").
		Or("description ILIKE ?", "%"+keyword+"%")

	// Get total count
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	err = query.Preload("Images").Preload("User").Offset(offset).Limit(limit).Order("created_at DESC").Find(&items).Error

	return items, count, err
}
