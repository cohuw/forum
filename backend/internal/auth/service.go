package auth

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"

	"forum_backend/internal/database"
	"forum_backend/internal/mail"
	"forum_backend/internal/models"
	"forum_backend/pkg/utils"
)

var (
	ErrEmailTaken    = errors.New("ерр email занят")
	ErrUsernameTaken = errors.New("ерр юзернейм занят")
	ErrInvalidCreds  = errors.New("ерр логин мимо")
	ErrCodeExpired   = errors.New("ерр код протух")
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type VerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

type ResendRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type PublicUser struct {
	ID         uint   `json:"id"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	IsVerified bool   `json:"is_verified"`
}

type TokenResponse struct {
	AccessToken string     `json:"access_token"`
	User        PublicUser `json:"user"`
}

func generateCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1_000_000))
}

func Register(req RegisterRequest) error {
	var existing models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		return ErrEmailTaken
	}
	if err := database.DB.Where("username = ?", req.Username).First(&existing).Error; err == nil {
		return ErrUsernameTaken
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("ерр хеш: %w", err)
	}

	code := generateCode()

	return database.DB.Transaction(func(tx *gorm.DB) error {
		user := models.User{
			Email:        req.Email,
			Username:     req.Username,
			PasswordHash: hash,
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		vc := models.VerificationCode{
			UserID:    user.ID,
			Code:      code,
			ExpiresAt: time.Now().Add(30 * time.Minute),
		}
		if err := tx.Create(&vc).Error; err != nil {
			return err
		}
		if err := mail.SendVerification(req.Email, code); err != nil {
			return fmt.Errorf("ерр почта: %w", err)
		}
		return nil
	})
}

func Login(req LoginRequest) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCreds
		}
		return nil, err
	}
	if !utils.CheckPassword(user.PasswordHash, req.Password) {
		return nil, ErrInvalidCreds
	}
	return &user, nil
}

func Verify(req VerifyRequest) (*models.User, error) {
	var vc models.VerificationCode
	err := database.DB.
		Joins("JOIN users ON users.id = verification_codes.user_id").
		Where("users.email = ? AND verification_codes.code = ? AND verification_codes.used = false AND verification_codes.expires_at > ?",
			req.Email, req.Code, time.Now()).
		First(&vc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCodeExpired
		}
		return nil, err
	}

	var user models.User
	return &user, database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vc).Update("used", true).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.User{}).Where("id = ?", vc.UserID).Update("is_verified", true).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", vc.UserID).First(&user).Error
	})
}

func Resend(req ResendRequest) error {
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil
	}
	if user.IsVerified {
		return nil
	}

	database.DB.Model(&models.VerificationCode{}).Where("user_id = ?", user.ID).Update("used", true)

	code := generateCode()
	vc := models.VerificationCode{
		UserID:    user.ID,
		Code:      code,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}
	if err := database.DB.Create(&vc).Error; err != nil {
		return err
	}
	return mail.SendVerification(req.Email, code)
}
