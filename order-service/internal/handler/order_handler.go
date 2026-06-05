package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tv-anagha/ecommerce-backend/order-service/internal/client"
	"github.com/tv-anagha/ecommerce-backend/order-service/internal/repository"
	"github.com/tv-anagha/ecommerce-backend/order-service/internal/service"
)

type OrderHandler struct {
	svc *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "order-service"})
}

type placeOrderBody struct {
	UserID uint `json:"userId" binding:"required"`
}

func (h *OrderHandler) PlaceOrder(c *gin.Context) {
	var body placeOrderBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	order, err := h.svc.PlaceOrder(body.UserID)
	if errors.Is(err, service.ErrEmptyCart) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if errors.Is(err, client.ErrProductNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cart contains a product that no longer exists"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": h.svc.CheckoutErrorMessage(err)})
		return
	}
	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	order, err := h.svc.GetOrder(uint(id))
	if errors.Is(err, repository.ErrOrderNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}
