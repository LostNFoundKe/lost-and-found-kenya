package repository

import (
	"gorm.io/gorm"
	"lostnfound-api/internal/models"
)

// ImageRepository handles database operations for images
type ImageRepository struct {
	db *gorm.DB
}

// NewImageRepository creates a new ImageRepository
func NewImageRepository(db *gorm.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

// Create adds a new image to the database
func (r *ImageRepository) Create(image *models.Image) error {
	return r.db.Create(image).Error
}

// GetByID retrieves an image by ID
func (r *ImageRepository) GetByID(id uint) (*models.Image, error) {
	var image models.Image
	err := r.db.First(&image, id).Error
	return &image, err
}

// GetByItemID retrieves all images for an item
func (r *ImageRepository) GetByItemID(itemID uint) ([]models.Image, error) {
	var images []models.Image
	err := r.db.Where("item_id = ?", itemID).Find(&images).Error
	return images, err
}

// Delete removes an image from the database
func (r *ImageRepository) Delete(id uint) error {
	return r.db.Delete(&models.Image{}, id).Error
}

// DeleteByItemID removes all images for an item
func (r *ImageRepository) DeleteByItemID(itemID uint) error {
	return r.db.Where("item_id = ?", itemID).Delete(&models.Image{}).Error
}
