package auth

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"forum_backend/internal/database"
	"forum_backend/internal/models"
	"forum_backend/pkg/utils"
)

type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ошибка": "ерр №3 рег парс"})
		return
	}
	if err := Register(req); err != nil {
		switch {
		case errors.Is(err, ErrEmailTaken):
			c.JSON(http.StatusConflict, gin.H{"ошибка": "ерр №4 email занят"})
		case errors.Is(err, ErrUsernameTaken):
			c.JSON(http.StatusConflict, gin.H{"ошибка": "ерр №5 юзернейм занят"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"ошибка": "ерр №6 рег фейл"})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"сообщение": "код улетел на почту"})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ошибка": "ерр №7 логин парс"})
		return
	}
	user, err := Login(req)
	if err != nil {
		if errors.Is(err, ErrInvalidCreds) {
			c.JSON(http.StatusUnauthorized, gin.H{"ошибка": "ерр №8 логин мимо"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"ошибка": "ерр №9 логин фейл"})
		return
	}
	issueTokens(c, user)
}

func (h *Handler) Verify(c *gin.Context) {
	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ошибка": "ерр №10 вериф парс"})
		return
	}
	user, err := Verify(req)
	if err != nil {
		if errors.Is(err, ErrCodeExpired) {
			c.JSON(http.StatusBadRequest, gin.H{"ошибка": "ерр №11 код протух"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"ошибка": "ерр №12 вериф фейл"})
		return
	}
	issueTokens(c, user)
}

func (h *Handler) Resend(c *gin.Context) {
	var req ResendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ошибка": "ерр №13 ресенд парс"})
		return
	}
	_ = Resend(req)
	c.JSON(http.StatusOK, gin.H{"сообщение": "если акк есть, код ушел"})
}

func (h *Handler) Refresh(c *gin.Context) {
	cookieVal, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"ошибка": "ерр №14 куки нет", "код": "нет_рефреша", "code": "token_expired"})
		return
	}
	claims, err := utils.ParseToken(cookieVal)
	if err != nil || claims.TokenType != "refresh" {
		c.JSON(http.StatusUnauthorized, gin.H{"ошибка": "ерр №15 токен стух", "код": "рефреш_истек", "code": "token_expired"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, claims.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"ошибка": "ерр №16 юзера нет"})
		return
	}
	issueTokens(c, &user)
}

func (h *Handler) Logout(c *gin.Context) {
	setRefreshCookie(c, "", -1)
	c.JSON(http.StatusOK, gin.H{"сообщение": "вышел"})
}

func issueTokens(c *gin.Context, user *models.User) {
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, user.Username, user.IsVerified)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ошибка": "ерр №1 токен"})
		return
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ошибка": "ерр №2 токен"})
		return
	}

	ttl, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_TTL_HOURS"))
	if ttl == 0 {
		ttl = 24
	}
	setRefreshCookie(c, refreshToken, ttl*3600)

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken: accessToken,
		User: PublicUser{
			ID:         user.ID,
			Email:      user.Email,
			Username:   user.Username,
			IsVerified: user.IsVerified,
		},
	})
}

func setRefreshCookie(c *gin.Context, value string, maxAgeSecs int) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("refresh_token", value, maxAgeSecs, "/", "", false, true)
}
