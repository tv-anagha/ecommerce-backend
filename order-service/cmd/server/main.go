package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tv-anagha/ecommerce-backend/order-service/internal/client"
	"github.com/tv-anagha/ecommerce-backend/order-service/internal/database"
	"github.com/tv-anagha/ecommerce-backend/order-service/internal/handler"
	"github.com/tv-anagha/ecommerce-backend/order-service/internal/repository"
	"github.com/tv-anagha/ecommerce-backend/order-service/internal/service"
)

func listenAddr() string {
	if port := os.Getenv("PORT"); port != "" {
		return "0.0.0.0:" + port
	}
	return "0.0.0.0:8084"
}

func main() {
	addr := listenAddr()
	fmt.Fprintf(os.Stderr, "order-service starting on %s\n", addr)

	database.Connect()

	repo := repository.NewOrderRepository()
	cartClient := client.NewCartClient()
	productClient := client.NewProductClient()
	svc := service.NewOrderService(repo, cartClient, productClient)
	h := handler.NewOrderHandler(svc)

	r := gin.Default()
	r.GET("/health", h.Health)
	r.POST("/orders", h.PlaceOrder)
	r.GET("/orders/:id", h.GetOrder)

	fmt.Fprintf(os.Stderr, "order-service listening on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("server failed:", err)
	}
}
