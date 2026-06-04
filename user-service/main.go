package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tv-anagha/ecommerce-backend/user-service/internal/database"
)

const listenAddr = "0.0.0.0:8080"

type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name" gorm:"column:product_name"`
	Price    float64 `json:"price" gorm:"column:price"`
	Category string `json:"category" gorm:"column:category"`
}

func main() {
	// Always print to stderr (visible even when GIN_MODE=release hides Gin debug).
	fmt.Fprintf(os.Stderr, "user-service starting on %s\n", listenAddr)

	database.Connect()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:5173",
		},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/ping-remote", func(c *gin.Context) {
		otherURL := os.Getenv("http://localhost:8080") // e.g. http://localhost:8081
		resp, err := http.Get(otherURL + "/ping")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil {
				log.Printf("close response body: %v", closeErr)
			}
		}()
	})

	r.GET("/products", func(c *gin.Context) {
		var products []Product
		if err := database.DB.Find(&products).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, products)
	})

	fmt.Fprintf(os.Stderr, "user-service listening on %s\n", listenAddr)
	if err := r.Run(listenAddr); err != nil {
		log.Fatal("server failed:", err)
	}
}
