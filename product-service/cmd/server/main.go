package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tv-anagha/ecommerce-backend/product-service/internal/database"
	"github.com/tv-anagha/ecommerce-backend/product-service/internal/handler"
	"github.com/tv-anagha/ecommerce-backend/product-service/internal/repository"
	"github.com/tv-anagha/ecommerce-backend/product-service/internal/service"
)

func listenAddr() string {
	if port := os.Getenv("PORT"); port != "" {
		return "0.0.0.0:" + port
	}
	return "0.0.0.0:8081"
}

func main() {
	addr := listenAddr()
	fmt.Fprintf(os.Stderr, "product-service starting on %s\n", addr)

	database.Connect()

	//Repository talks directly to the database.
	repo := repository.NewProductRepository()
	//Service talks to the repository. Service contains business logic.
	svc := service.NewProductService(repo)
	//Handlers receive HTTP requests and pass them to the service.
	h := handler.NewProductHandler(svc)


	// Route registration
	r := gin.Default()
	// Gateway registers Gin routes and mounts the reverse proxy:

	r.GET("/health", h.Health)
	r.GET("/products", h.ListProducts)
	r.GET("/products/:id", h.GetProduct)

	fmt.Fprintf(os.Stderr, "product-service listening on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("server failed:", err)
	}
}
