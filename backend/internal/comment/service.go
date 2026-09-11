package comment

import (
	"errors"
	"fmt"
	"html"
	"time"

	"forum_backend/internal/database"
	"forum_backend/internal/models"
)

var ErrForbidden = errors.New("ерр не твой коммент")

type CreateRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

type CommentResponse struct {
	ID        uint      `json:"id"`
	PostID    uint      `json:"post_id"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Likes     int64     `json:"likes"`
	Dislikes  int64     `json:"dislikes"`
	UserVote  int8      `json:"user_vote"`
}

func GetByPost(postID uint, currentUserID *uint) ([]CommentResponse, error) {
	var rows []CommentResponse
	err := database.DB.Model(&models.Comment{}).
		Select(`comments.id, comments.post_id, comments.user_id, users.username,
			comments.content, comments.created_at,
			COALESCE(SUM(CASE WHEN v.value = 1 THEN 1 ELSE 0 END), 0) AS likes,
			COALESCE(SUM(CASE WHEN v.value = -1 THEN 1 ELSE 0 END), 0) AS dislikes`).
		Joins("JOIN users ON users.id = comments.user_id").
		Joins("LEFT JOIN votes v ON v.target_type = 'comment' AND v.target_id = comments.id").
		Where("comments.post_id = ?", postID).
		Group("comments.id, users.username").
		Order("comments.created_at ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	if currentUserID != nil {
		for i := range rows {
			var v models.Vote
			if err := database.DB.Where("user_id = ? AND target_type = 'comment' AND target_id = ?",
				*currentUserID, rows[i].ID).First(&v).Error; err == nil {
				rows[i].UserVote = v.Value
			}
		}
	}
	if rows == nil {
		rows = []CommentResponse{}
	}
	return rows, nil
}

func Create(postID, userID uint, req CreateRequest) (*CommentResponse, error) {
	c := models.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: html.EscapeString(req.Content),
	}
	if err := database.DB.Create(&c).Error; err != nil {
		return nil, fmt.Errorf("ерр криейт коммент: %w", err)
	}
	var resp CommentResponse
	err := database.DB.Model(&models.Comment{}).
		Select("comments.id, comments.post_id, comments.user_id, users.username, comments.content, comments.created_at").
		Joins("JOIN users ON users.id = comments.user_id").
		Where("comments.id = ?", c.ID).
		Scan(&resp).Error
	return &resp, err
}

func Delete(commentID, userID uint) error {
	result := database.DB.Where("id = ? AND user_id = ?", commentID, userID).Delete(&models.Comment{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrForbidden
	}
	return nil
}
