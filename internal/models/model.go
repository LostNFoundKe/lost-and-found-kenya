package models

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Model struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate hook to generate a UUID before saving a new record
func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New() // Generates a new unique UUID
	return nil
}

// JsonResponse defines the API response structure
type JsonResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func ResponseJson(c *gin.Context, status int, message string, data any) {
	response := JsonResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
	c.JSON(status, response)

}

func ParseIntOrDefault(val string, defaultVal int) int {
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return intVal
}
