package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"forum_backend/internal/auth"
	"forum_backend/internal/category"
	"forum_backend/internal/comment"
	"forum_backend/internal/database"
	"forum_backend/internal/middleware"
	"forum_backend/internal/post"
	"forum_backend/internal/vote"
)

func main() {
	if err := godotenv.Load("configs/.env"); err != nil {
		log.Println(".енв нет")
	}

	if err := database.Connect(); err != nil {
		log.Fatalf("ерр бд: %v", err)
	}
	log.Println("бд коннект")

	if err := database.Migrate(); err != nil {
		log.Fatalf("ерр миграции: %v", err)
	}
	log.Println("миграции успешно")

	authH := auth.NewHandler()
	categoryH := category.NewHandler()
	postH := post.NewHandler()
	commentH := comment.NewHandler()
	voteH := vote.NewHandler()

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	api := r.Group("/api")
	{
		authG := api.Group("/auth")
		{
			authG.POST("/register", authH.Register)
			authG.POST("/login", authH.Login)
			authG.POST("/verify", authH.Verify)
			authG.POST("/resend", authH.Resend)
			authG.GET("/refresh", authH.Refresh)
			authG.POST("/logout", authH.Logout)
		}

		api.GET("/categories", categoryH.List)

		api.GET("/posts", middleware.OptionalAuth(), postH.List)
		api.POST("/posts", middleware.RequireVerified(), postH.Create)
		api.GET("/posts/:id", middleware.OptionalAuth(), postH.Get)
		api.DELETE("/posts/:id", middleware.Auth(), postH.Delete)

		api.GET("/posts/:id/comments", middleware.OptionalAuth(), commentH.ListByPost)
		api.POST("/posts/:id/comments", middleware.RequireVerified(), commentH.Create)
		api.DELETE("/comments/:id", middleware.Auth(), commentH.Delete)

		api.POST("/votes", middleware.RequireVerified(), voteH.Vote)
	}

	frontendDir := os.Getenv("FRONTEND_DIR")
	if frontendDir == "" {
		frontendDir = "../frontend/dist"
	}
	if _, err := os.Stat(frontendDir); err == nil {
		r.Static("/assets", frontendDir+"/assets")
		r.StaticFile("/", frontendDir+"/index.html")
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if len(path) < 4 || path[:4] != "/api" {
				c.File(frontendDir + "/index.html")
			}
		})
		log.Printf("фронт из: %s", frontendDir)
	} else {
		log.Printf("фронт %q нет / билд", frontendDir)
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("серв ран: http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("ерр серв: %v", err)
	}
}