package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tv-anagha/ecommerce-backend/order-service/internal/client"
	"github.com/tv-anagha/ecommerce-backend/order-service/internal/database"
	"github.com/tv-anagha/ecommerce-backend/order-service/internal/events"
	"github.com/tv-anagha/ecommerce-backend/order-service/internal/handler"
	kafkapkg "github.com/tv-anagha/ecommerce-backend/order-service/internal/kafka"
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

	//Creates repository for the order model.
	repo := repository.NewOrderRepository()
	//Creates HTTP client for Cart Service to get cart details.
	cartClient := client.NewCartClient()
	//Creates HTTP client for Product Service to get product details.
	productClient := client.NewProductClient()
	//Creates connection to Kafka broker.
	producer := kafkapkg.NewProducer()
	//Closes the connection to Kafka broker.
	defer func() { _ = producer.Close() }()
//Creates publisher for the order placed event and wraps the Kafka producer.		
	publisher := events.NewPublisher(producer)
	//Creates service for the order model.
	svc := service.NewOrderService(repo, cartClient, productClient, publisher)
	//Creates handler for the order model.
	h := handler.NewOrderHandler(svc)

	//Creates HTTP server.
	r := gin.Default()
	r.GET("/health", h.Health)

	r.POST("/orders", h.PlaceOrder)
	r.GET("/orders/:id", h.GetOrder)

	fmt.Fprintf(os.Stderr, "order-service listening on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("server failed:", err)
	}
}
