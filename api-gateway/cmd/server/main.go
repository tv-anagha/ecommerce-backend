package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tv-anagha/ecommerce-backend/api-gateway/internal/proxy"
)

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func corsOrigins() []string {
	raw := env("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173,http://127.0.0.1:3000,http://127.0.0.1:5173")
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, o := range parts {
		o = strings.TrimSpace(o)
		if o != "" {
			origins = append(origins, o)
		}
	}
	return origins
}

func listenAddr() string {
	if port := os.Getenv("PORT"); port != "" {
		return "0.0.0.0:" + port
	}
	return "0.0.0.0:8080"
}

func mountProxy(r *gin.Engine, gatewayPrefix, rewritePrefix, targetURL string) {
	target := proxy.MustParseURL(targetURL)
	h := proxy.Handler(target, gatewayPrefix, rewritePrefix)
	r.Any(gatewayPrefix+"/*path", gin.WrapH(h))
	r.Any(gatewayPrefix, gin.WrapH(h))
}

func checkHealth(url string) gin.H {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return gin.H{"status": "down", "error": err.Error()}
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return gin.H{"status": "down", "httpStatus": resp.StatusCode}
	}
	return gin.H{"status": "ok"}
}

func main() {
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")

	addr := listenAddr()
	fmt.Fprintf(os.Stderr, "api-gateway starting on %s\n", addr)

	productURL := env("PRODUCT_SERVICE_URL", "http://localhost:8081")

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: corsOrigins(),
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "api-gateway",
			"services": gin.H{
				"product": checkHealth(productURL + "/health"),
			},
		})
	})

	mountProxy(r, "/api/products", "/products", productURL)

	fmt.Fprintf(os.Stderr, "api-gateway listening on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("server failed:", err)
	}
}
