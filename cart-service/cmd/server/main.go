package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tv-anagha/ecommerce-backend/cart-service/internal/client"
	"github.com/tv-anagha/ecommerce-backend/cart-service/internal/database"
	"github.com/tv-anagha/ecommerce-backend/cart-service/internal/handler"
	"github.com/tv-anagha/ecommerce-backend/cart-service/internal/repository"
	"github.com/tv-anagha/ecommerce-backend/cart-service/internal/service"
)

func listenAddr() string {
	if port := os.Getenv("PORT"); port != "" {
		return "0.0.0.0:" + port
	}
	return "0.0.0.0:8083"
}

func main() {
	addr := listenAddr()
	fmt.Fprintf(os.Stderr, "cart-service starting on %s\n", addr)

	database.Connect()

	repo := repository.NewCartRepository()
	productClient := client.NewProductClient()
	svc := service.NewCartService(repo, productClient)
	h := handler.NewCartHandler(svc)

	r := gin.Default()
	r.GET("/health", h.Health)
	r.GET("/carts/:userId", h.GetCart)
	r.POST("/carts/:userId/items", h.AddItem)
	r.PATCH("/carts/:userId/items/:productId", h.UpdateItem)
	r.DELETE("/carts/:userId/items/:productId", h.RemoveItem)

	fmt.Fprintf(os.Stderr, "cart-service listening on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("server failed:", err)
	}
}
