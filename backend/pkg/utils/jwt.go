package utils

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID     uint   `json:"user_id"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	IsVerified bool   `json:"is_verified"`
	TokenType  string `json:"token_type"`
	jwt.RegisteredClaims
}

func getSecret() []byte { return []byte(os.Getenv("JWT_SECRET")) }

func GenerateAccessToken(userID uint, email, username string, isVerified bool) (string, error) {
	ttl, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_TTL_MINUTES"))
	if ttl == 0 {
		ttl = 20
	}
	claims := Claims{
		UserID: userID, Email: email, Username: username,
		IsVerified: isVerified, TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(ttl) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(getSecret())
}

func GenerateRefreshToken(userID uint) (string, error) {
	ttl, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_TTL_HOURS"))
	if ttl == 0 {
		ttl = 24
	}
	claims := Claims{
		UserID: userID, TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(ttl) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(getSecret())
}

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("ерр метод подписи")
		}
		return getSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("ерр токен невалид")
	}
	return claims, nil
}
