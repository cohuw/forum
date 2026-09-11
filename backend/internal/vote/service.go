package vote

import (
	"errors"
	"fmt"

	"forum_backend/internal/database"
	"forum_backend/internal/models"
)

var (
	ErrInvalidTarget = errors.New("ерр таргет воут")
	ErrInvalidValue  = errors.New("ерр значение воут")
)

type VoteRequest struct {
	TargetType string `json:"target_type" binding:"required"`
	TargetID   uint   `json:"target_id" binding:"required"`
	Value      int8   `json:"value"`
}

type VoteResponse struct {
	Likes    int64 `json:"likes"`
	Dislikes int64 `json:"dislikes"`
	UserVote int8  `json:"user_vote"`
}

func Vote(userID uint, req VoteRequest) (*VoteResponse, error) {
	if req.TargetType != "post" && req.TargetType != "comment" {
		return nil, ErrInvalidTarget
	}
	if req.Value != 1 && req.Value != -1 && req.Value != 0 {
		return nil, ErrInvalidValue
	}

	if req.Value == 0 {
		err := database.DB.
			Where("user_id = ? AND target_type = ? AND target_id = ?", userID, req.TargetType, req.TargetID).
			Delete(&models.Vote{}).Error
		if err != nil {
			return nil, fmt.Errorf("ерр удаление воут: %w", err)
		}
	} else {
		err := database.DB.Exec(`
			INSERT INTO votes (user_id, target_type, target_id, value)
			VALUES (?, ?, ?, ?)
			ON CONFLICT (user_id, target_type, target_id)
			DO UPDATE SET value = EXCLUDED.value`,
			userID, req.TargetType, req.TargetID, req.Value,
		).Error
		if err != nil {
			return nil, fmt.Errorf("ерр сохранение воут: %w", err)
		}
	}
	return getCounts(userID, req.TargetType, req.TargetID)
}

func getCounts(userID uint, targetType string, targetID uint) (*VoteResponse, error) {
	type agg struct {
		Likes    int64
		Dislikes int64
	}
	var r agg
	database.DB.Model(&models.Vote{}).
		Select(`COALESCE(SUM(CASE WHEN value = 1 THEN 1 ELSE 0 END), 0) AS likes,
			COALESCE(SUM(CASE WHEN value = -1 THEN 1 ELSE 0 END), 0) AS dislikes`).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Scan(&r)

	var userVote int8
	var v models.Vote
	if err := database.DB.Where("user_id = ? AND target_type = ? AND target_id = ?",
		userID, targetType, targetID).First(&v).Error; err == nil {
		userVote = v.Value
	}
	return &VoteResponse{Likes: r.Likes, Dislikes: r.Dislikes, UserVote: userVote}, nil
}
