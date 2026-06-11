package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tv-anagha/ecommerce-backend/fulfillment-service/internal/consumer"
	"github.com/tv-anagha/ecommerce-backend/fulfillment-service/internal/handler"
)

func listenAddr() string {
	if port := os.Getenv("PORT"); port != "" {
		return "0.0.0.0:" + port
	}
	return "0.0.0.0:8086"
}

func main() {
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")

	addr := listenAddr()
	fmt.Fprintf(os.Stderr, "fulfillment-service starting on %s\n", addr)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	orderConsumer := consumer.NewOrderConsumer()
	defer func() { _ = orderConsumer.Close() }()

	// The Kafka consumer runs in its own goroutine so it can block on ReadMessage
	// without stalling the HTTP server. /health stays responsive for Docker and
	// orchestration probes while the consumer waits for the next order.placed event.
	go func() {
		if err := orderConsumer.Run(ctx); err != nil {
			log.Printf("kafka consumer stopped: %v", err)
		}
	}()

	h := handler.NewHealthHandler()
	r := gin.Default()
	r.GET("/health", h.Health)

	server := &http.Server{Addr: addr, Handler: r}
	go func() {
		fmt.Fprintf(os.Stderr, "fulfillment-service listening on %s\n", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server failed:", err)
		}
	}()

	<-ctx.Done()
	_ = server.Shutdown(context.Background())
}
