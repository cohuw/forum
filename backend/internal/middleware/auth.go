package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"forum_backend/pkg/utils"
)

const ClaimsKey = "claims"

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"ошибка": "ерр №40 нет хедера",
				"код":    "нет_токена",
				"code":   "no_token",
			})
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"ошибка": "ерр №41 токен сдох",
				"код":    "токен_истек",
				"code":   "token_expired",
			})
			return
		}
		if claims.TokenType != "access" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"ошибка": "ерр №42 тип токена"})
			return
		}
		c.Set(ClaimsKey, claims)
		c.Next()
	}
}

func RequireVerified() gin.HandlerFunc {
	return func(c *gin.Context) {
		Auth()(c)
		if c.IsAborted() {
			return
		}
		claims := GetClaims(c)
		if !claims.IsVerified {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"ошибка": "ерр №43 не подтвержден"})
			return
		}
		c.Next()
	}
}

func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if strings.HasPrefix(header, "Bearer ") {
			tokenStr := strings.TrimPrefix(header, "Bearer ")
			if claims, err := utils.ParseToken(tokenStr); err == nil && claims.TokenType == "access" {
				c.Set(ClaimsKey, claims)
			}
		}
		c.Next()
	}
}

func GetClaims(c *gin.Context) *utils.Claims {
	val, exists := c.Get(ClaimsKey)
	if !exists {
		return nil
	}
	claims, _ := val.(*utils.Claims)
	return claims
}
