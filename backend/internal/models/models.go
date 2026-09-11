package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	Username     string    `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"not null" json:"-"`
	IsVerified   bool      `gorm:"default:false" json:"is_verified"`
}

type VerificationCode struct {
	gorm.Model
	UserID    uint      `gorm:"not null;index"`
	Code      string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	Used      bool      `gorm:"default:false"`
}

type Category struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `gorm:"uniqueIndex;not null" json:"name"`
	Slug  string `gorm:"uniqueIndex;not null" json:"slug"`
	Posts []Post `gorm:"many2many:post_categories;" json:"-"`
}

type Post struct {
	gorm.Model
	UserID     uint       `gorm:"not null;index"`
	User       User       `gorm:"foreignKey:UserID"`
	Title      string     `gorm:"not null"`
	Content    string     `gorm:"not null;type:text"`
	Categories []Category `gorm:"many2many:post_categories;"`
	Comments   []Comment  `gorm:"foreignKey:PostID"`
}

type Comment struct {
	gorm.Model
	PostID  uint   `gorm:"not null;index"`
	UserID  uint   `gorm:"not null;index"`
	User    User   `gorm:"foreignKey:UserID"`
	Content string `gorm:"not null;type:text"`
}

type Vote struct {
	ID         uint   `gorm:"primaryKey"`
	UserID     uint   `gorm:"not null;uniqueIndex:idx_vote_unique"`
	TargetType string `gorm:"not null;uniqueIndex:idx_vote_unique"`
	TargetID   uint   `gorm:"not null;uniqueIndex:idx_vote_unique"`
	Value      int8   `gorm:"not null"`
}
