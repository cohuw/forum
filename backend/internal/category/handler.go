package category

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"forum_backend/internal/models"
)

type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

func (h *Handler) List(c *gin.Context) {
	cats, err := GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ошибка": "ерр №17 кат"})
		return
	}
	if cats == nil {
		cats = []models.Category{}
	}
	c.JSON(http.StatusOK, cats)
}
