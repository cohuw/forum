package post

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"forum_backend/internal/middleware"
)

type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

func (h *Handler) List(c *gin.Context) {
	filter := Filter{}
	claims := middleware.GetClaims(c)

	if catStr := c.Query("category"); catStr != "" {
		if id, err := strconv.ParseUint(catStr, 10, 64); err == nil {
			uid := uint(id)
			filter.CategoryID = &uid
		}
	}
	if c.Query("my") == "true" {
		if claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"ошибка": "ерр №18 мои посты без входа"})
			return
		}
		filter.UserID = &claims.UserID
	}
	if c.Query("liked") == "true" {
		if claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"ошибка": "ерр №19 лайки без входа"})
			return
		}
		filter.LikedByID = &claims.UserID
	}

	var currentUserID *uint
	if claims != nil {
		currentUserID = &claims.UserID
	}

	posts, err := GetPosts(filter, currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ошибка": "ерр №20 посты"})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ошибка": "ерр №21 ид поста"})
		return
	}
	claims := middleware.GetClaims(c)
	var currentUserID *uint
	if claims != nil {
		currentUserID = &claims.UserID
	}
	p, err := GetByID(uint(id), currentUserID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"ошибка": "ерр №22 пост 404"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"ошибка": "ерр №23 гет пост"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *Handler) Create(c *gin.Context) {
	claims := middleware.GetClaims(c)
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ошибка": "ерр №24 пост парс"})
		return
	}
	p, err := Create(claims.UserID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ошибка": "ерр №25 пост криейт"})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *Handler) Delete(c *gin.Context) {
	claims := middleware.GetClaims(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ошибка": "ерр №26 ид делит"})
		return
	}
	if err := Delete(uint(id), claims.UserID); err != nil {
		if errors.Is(err, ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"ошибка": "ерр №27 не твой пост"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"ошибка": "ерр №28 делит пост"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"сообщение": "пост удален"})
}
