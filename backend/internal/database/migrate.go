package database

import (
	"log"

	"forum_backend/internal/models"
)

func Migrate() error {
	err := DB.AutoMigrate(
		&models.User{},
		&models.VerificationCode{},
		&models.Category{},
		&models.Post{},
		&models.Comment{},
		&models.Vote{},
	)
	if err != nil {
		return err
	}

	var count int64
	DB.Model(&models.Category{}).Count(&count)
	if count == 0 {
		seed := []models.Category{
			{Name: "Общее", Slug: "general"},
			{Name: "Технологии", Slug: "technology"},
			{Name: "Наука", Slug: "science"},
			{Name: "Оффтоп", Slug: "off-topic"},
			{Name: "Объявления", Slug: "announcements"},
		}
		if err := DB.Create(&seed).Error; err != nil {
			log.Printf("ерр сид категорий: %v", err)
		} else {
			log.Println("категории ок")
		}
	}
	return nil
}
