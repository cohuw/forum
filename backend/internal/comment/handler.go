package comment

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"forum_backend/internal/middleware"
)

type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

func (h *Handler) ListByPost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ошибка": "ерр №29 ид коммент"})
		return
	}
	claims := middleware.GetClaims(c)
	var currentUserID *uint
	if claims != nil {
		currentUserID = &claims.UserID
	}
	comments, err := GetByPost(uint(postID), currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ошибка": "ерр №30 гет комменты"})
		return
	}
	c.JSON(http.StatusOK, comments)
}

func (h *Handler) Create(c *gin.Context) {
	claims := middleware.GetClaims(c)
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ошибка": "ерр №31 ид коммент пост"})
		return
	}
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ошибка": "ерр №32 коммент парс"})
		return
	}
	comment, err := Create(uint(postID), claims.UserID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ошибка": "ерр №33 коммент криейт"})
		return
	}
	c.JSON(http.StatusCreated, comment)
}

func (h *Handler) Delete(c *gin.Context) {
	claims := middleware.GetClaims(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ошибка": "ерр №34 ид делит коммент"})
		return
	}
	if err := Delete(uint(id), claims.UserID); err != nil {
		if errors.Is(err, ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"ошибка": "ерр №35 не твой коммент"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"ошибка": "ерр №36 делит коммент"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"сообщение": "коммент удален"})
}
