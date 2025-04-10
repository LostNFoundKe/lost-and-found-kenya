package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"lostnfound-api/internal/models"
	"lostnfound-api/internal/service"
	"net/http"
)

// ItemHandler handles HTTP requests for items
type ItemHandler struct {
	service *service.ItemService
}

// NewItemHandler creates a new ItemHandler
func NewItemHandler(service *service.ItemService) *ItemHandler {
	return &ItemHandler{service: service}
}

// Create handles the creation of a new item
func (h *ItemHandler) Create(c *gin.Context) {
	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		models.ResponseJson(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		models.ResponseJson(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	item.UserID = userID.(uuid.UUID)

	if err := h.service.Create(&item); err != nil {
		models.ResponseJson(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	models.ResponseJson(c, http.StatusCreated, "Item created successfully", item)
}

// GetByID handles retrieval of an item by ID
func (h *ItemHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		models.ResponseJson(c, http.StatusBadRequest, "invalid ID", nil)
		return
	}

	item, err := h.service.GetByID(id)
	if err != nil {
		models.ResponseJson(c, http.StatusNotFound, "item not found", nil)
		return
	}

	models.ResponseJson(c, http.StatusOK, "Item retrieved successfully", item)
}

// List handles retrieval of items with filtering
func (h *ItemHandler) List(c *gin.Context) {
	status := c.Query("status")
	category := c.Query("category")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	// Convert page and limit to integers
	pageInt := 1
	limitInt := 10

	// We don't need to check for errors here since we have default values
	pageInt = models.ParseIntOrDefault(page, 1)
	limitInt = models.ParseIntOrDefault(limit, 10)

	items, count, err := h.service.List(status, category, pageInt, limitInt)
	if err != nil {
		models.ResponseJson(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	responseData := gin.H{
		"items": items,
		"total": count,
		"page":  pageInt,
		"limit": limitInt,
	}

	models.ResponseJson(c, http.StatusOK, "Items retrieved successfully", responseData)
}

// Update handles updating an existing item
func (h *ItemHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		models.ResponseJson(c, http.StatusBadRequest, "invalid ID", nil)
		return
	}

	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		models.ResponseJson(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	item.ID = id

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		models.ResponseJson(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	// Get the existing item to check ownership
	existingItem, err := h.service.GetByID(id)
	if err != nil {
		models.ResponseJson(c, http.StatusNotFound, "item not found", nil)
		return
	}

	// Check if the user owns the item or is an admin
	isAdmin, _ := c.Get("isAdmin")
	if existingItem.UserID != userID.(uuid.UUID) && isAdmin != true {
		models.ResponseJson(c, http.StatusForbidden, "not authorized to update this item", nil)
		return
	}

	if err := h.service.Update(&item); err != nil {
		models.ResponseJson(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	models.ResponseJson(c, http.StatusOK, "Item updated successfully", item)
}

// Delete handles removal of an item
func (h *ItemHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		models.ResponseJson(c, http.StatusBadRequest, "invalid ID", nil)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		models.ResponseJson(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	// Get the existing item to check ownership
	existingItem, err := h.service.GetByID(id)
	if err != nil {
		models.ResponseJson(c, http.StatusNotFound, "item not found", nil)
		return
	}

	// Check if the user owns the item or is an admin
	isAdmin, _ := c.Get("isAdmin")
	if existingItem.UserID != userID.(uuid.UUID) && isAdmin != true {
		models.ResponseJson(c, http.StatusForbidden, "not authorized to delete this item", nil)
		return
	}

	if err := h.service.Delete(id); err != nil {
		models.ResponseJson(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	models.ResponseJson(c, http.StatusOK, "Item deleted successfully", nil)
}
