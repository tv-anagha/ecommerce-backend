package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tv-anagha/ecommerce-backend/user-service/internal/database"
	"github.com/tv-anagha/ecommerce-backend/user-service/internal/handler"
	"github.com/tv-anagha/ecommerce-backend/user-service/internal/repository"
	"github.com/tv-anagha/ecommerce-backend/user-service/internal/service"
)

func listenAddr() string {
	if port := os.Getenv("PORT"); port != "" {
		return "0.0.0.0:" + port
	}
	return "0.0.0.0:8082"
}

func corsOrigins() []string {
	raw := os.Getenv("CORS_ALLOWED_ORIGINS")
	if raw == "" {
		raw = "http://localhost:3000,http://localhost:5173,http://127.0.0.1:3000,http://127.0.0.1:5173"
	}
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, o := range parts {
		if o = strings.TrimSpace(o); o != "" {
			origins = append(origins, o)
		}
	}
	return origins
}

func main() {
	addr := listenAddr()
	fmt.Fprintf(os.Stderr, "user-service starting on %s\n", addr)

	database.Connect()

	repo := repository.NewUserRepository()
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)

	r := gin.Default()
	// Browsers call :8082 directly from the UI — CORS required (curl skips this).
	r.Use(cors.New(cors.Config{
		AllowOrigins: corsOrigins(),
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	r.GET("/health", h.Health)
	// Canonical paths
	r.POST("/users", h.Register)
	r.POST("/users/login", h.Login)
	r.GET("/users/me", h.Me)
	// Aliases many frontends expect (fixes UI 404 while curl used /users)
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)

	fmt.Fprintf(os.Stderr, "user-service listening on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("server failed:", err)
	}
}
