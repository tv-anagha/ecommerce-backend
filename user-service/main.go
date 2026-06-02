package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

const listenAddr = "0.0.0.0:8080"

func main() {
	// Always print to stderr (visible even when GIN_MODE=release hides Gin debug).
	fmt.Fprintf(os.Stderr, "user-service starting on %s\n", listenAddr)

	r := gin.Default()

    r.Use(cors.New(cors.Config{
        AllowOrigins: []string{"http://localhost:5173"},
        AllowMethods: []string{"GET"},
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
		defer resp.Body.Close()
	})

	fmt.Fprintf(os.Stderr, "user-service listening on %s\n", listenAddr)
	if err := r.Run(listenAddr); err != nil {
		log.Fatal("server failed:", err)
	}
}
