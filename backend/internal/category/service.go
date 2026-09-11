package category

import (
	"forum_backend/internal/database"
	"forum_backend/internal/models"
)

func GetAll() ([]models.Category, error) {
	var cats []models.Category
	err := database.DB.Order("name").Find(&cats).Error
	return cats, err
}