package post

import (
	"errors"
	"fmt"
	"html"
	"time"

	"forum_backend/internal/database"
	"forum_backend/internal/models"
)

var (
	ErrNotFound  = errors.New("ерр пост 404")
	ErrForbidden = errors.New("ерр не твой пост")
)

type CreateRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=200"`
	Content     string `json:"content" binding:"required,min=1"`
	CategoryIDs []uint `json:"category_ids"`
}

type PostResponse struct {
	ID           uint              `json:"id"`
	UserID       uint              `json:"user_id"`
	Username     string            `json:"username"`
	Title        string            `json:"title"`
	Content      string            `json:"content"`
	CreatedAt    time.Time         `json:"created_at"`
	Categories   []models.Category `gorm:"-" json:"categories"`
	Likes        int64             `json:"likes"`
	Dislikes     int64             `json:"dislikes"`
	UserVote     int8              `json:"user_vote"`
	CommentCount int64             `json:"comment_count"`
}

type Filter struct {
	CategoryID *uint
	UserID     *uint
	LikedByID  *uint
}

func GetPosts(filter Filter, currentUserID *uint) ([]PostResponse, error) {
	q := database.DB.Model(&models.Post{}).
		Select(`posts.id, posts.user_id, users.username, posts.title, posts.content, posts.created_at,
			COALESCE(SUM(CASE WHEN v.value = 1 THEN 1 ELSE 0 END), 0) AS likes,
			COALESCE(SUM(CASE WHEN v.value = -1 THEN 1 ELSE 0 END), 0) AS dislikes,
			COUNT(DISTINCT c.id) AS comment_count`).
		Joins("JOIN users ON users.id = posts.user_id").
		Joins("LEFT JOIN votes v ON v.target_type = 'post' AND v.target_id = posts.id").
		Joins("LEFT JOIN comments c ON c.post_id = posts.id").
		Group("posts.id, users.username").
		Order("posts.created_at DESC").
		Limit(50)

	if filter.CategoryID != nil {
		q = q.Joins("JOIN post_categories pc ON pc.post_id = posts.id AND pc.category_id = ?", *filter.CategoryID)
	}
	if filter.UserID != nil {
		q = q.Where("posts.user_id = ?", *filter.UserID)
	}
	if filter.LikedByID != nil {
		q = q.Joins("JOIN votes liked_v ON liked_v.target_type = 'post' AND liked_v.target_id = posts.id AND liked_v.user_id = ? AND liked_v.value = 1", *filter.LikedByID)
	}

	var rows []PostResponse
	if err := q.Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("ерр посты: %w", err)
	}

	for i := range rows {
		var cats []models.Category
		database.DB.Model(&models.Category{}).
			Joins("JOIN post_categories pc ON pc.category_id = categories.id").
			Where("pc.post_id = ?", rows[i].ID).
			Find(&cats)
		if cats == nil {
			cats = []models.Category{}
		}
		rows[i].Categories = cats

		if currentUserID != nil {
			var v models.Vote
			if err := database.DB.Where("user_id = ? AND target_type = 'post' AND target_id = ?",
				*currentUserID, rows[i].ID).First(&v).Error; err == nil {
				rows[i].UserVote = v.Value
			}
		}
	}
	if rows == nil {
		rows = []PostResponse{}
	}
	return rows, nil
}

func GetByID(id uint, currentUserID *uint) (*PostResponse, error) {
	var row PostResponse
	err := database.DB.Model(&models.Post{}).
		Select(`posts.id, posts.user_id, users.username, posts.title, posts.content, posts.created_at,
			COALESCE(SUM(CASE WHEN v.value = 1 THEN 1 ELSE 0 END), 0) AS likes,
			COALESCE(SUM(CASE WHEN v.value = -1 THEN 1 ELSE 0 END), 0) AS dislikes,
			COUNT(DISTINCT c.id) AS comment_count`).
		Joins("JOIN users ON users.id = posts.user_id").
		Joins("LEFT JOIN votes v ON v.target_type = 'post' AND v.target_id = posts.id").
		Joins("LEFT JOIN comments c ON c.post_id = posts.id").
		Where("posts.id = ?", id).
		Group("posts.id, users.username").
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, ErrNotFound
	}

	var cats []models.Category
	database.DB.Model(&models.Category{}).
		Joins("JOIN post_categories pc ON pc.category_id = categories.id").
		Where("pc.post_id = ?", id).
		Find(&cats)
	if cats == nil {
		cats = []models.Category{}
	}
	row.Categories = cats

	if currentUserID != nil {
		var v models.Vote
		if err := database.DB.Where("user_id = ? AND target_type = 'post' AND target_id = ?",
			*currentUserID, id).First(&v).Error; err == nil {
			row.UserVote = v.Value
		}
	}
	return &row, nil
}

func Create(userID uint, req CreateRequest) (*PostResponse, error) {
	post := models.Post{
		UserID:  userID,
		Title:   html.EscapeString(req.Title),
		Content: html.EscapeString(req.Content),
	}
	if len(req.CategoryIDs) > 0 {
		var cats []models.Category
		if err := database.DB.Where("id IN ?", req.CategoryIDs).Find(&cats).Error; err != nil {
			return nil, err
		}
		post.Categories = cats
	}
	if err := database.DB.Create(&post).Error; err != nil {
		return nil, fmt.Errorf("ерр криейт пост: %w", err)
	}
	return GetByID(post.ID, &userID)
}

func Delete(postID, userID uint) error {
	result := database.DB.Where("id = ? AND user_id = ?", postID, userID).Delete(&models.Post{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrForbidden
	}
	return nil
}
