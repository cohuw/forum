package vote

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"forum_backend/internal/middleware"
)

type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

func (h *Handler) Vote(c *gin.Context) {
	claims := middleware.GetClaims(c)
	var req VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ошибка": "ерр №37 воут парс"})
		return
	}
	resp, err := Vote(claims.UserID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidTarget), errors.Is(err, ErrInvalidValue):
			c.JSON(http.StatusBadRequest, gin.H{"ошибка": "ерр №38 кривой воут"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"ошибка": "ерр №39 воут"})
		}
		return
	}
	c.JSON(http.StatusOK, resp)
}
