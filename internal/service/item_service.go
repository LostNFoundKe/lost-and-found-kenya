package service

import (
	"errors"
	"github.com/google/uuid"
	"lostnfound-api/internal/models"
	"lostnfound-api/internal/repository"
)

// ItemService provides business logic for items
type ItemService struct {
	repo *repository.ItemRepository
}

// NewItemService creates a new ItemService
func NewItemService(repo *repository.ItemRepository) *ItemService {
	return &ItemService{repo: repo}
}

// Create adds a new item
func (s *ItemService) Create(item *models.Item) error {
	// Validate item
	if item.Title == "" {
		return errors.New("title is required")
	}

	return s.repo.Create(item)
}

// GetByID retrieves an item by ID
func (s *ItemService) GetByID(id uuid.UUID) (*models.Item, error) {
	return s.repo.GetByID(id)
}

// List retrieves items with filtering options
func (s *ItemService) List(status string, category string, page, limit int) ([]models.Item, int64, error) {
	// Default pagination values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	return s.repo.List(status, category, page, limit)
}

// Update updates an existing item
func (s *ItemService) Update(item *models.Item) error {
	// Validate item
	if item.Title == "" {
		return errors.New("title is required")
	}

	// Check if item exists
	_, err := s.repo.GetByID(item.ID)
	if err != nil {
		return errors.New("item not found")
	}

	return s.repo.Update(item)
}

// Delete removes an item
func (s *ItemService) Delete(id uuid.UUID) error {
	// Check if item exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("item not found")
	}

	return s.repo.Delete(id)
}

// Search searches items by keyword
func (s *ItemService) Search(keyword string, page, limit int) ([]models.Item, int64, error) {
	// Default pagination values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	return s.repo.SearchByKeyword(keyword, page, limit)
}
