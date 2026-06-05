package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tv-anagha/ecommerce-backend/cart-service/internal/client"
	"github.com/tv-anagha/ecommerce-backend/cart-service/internal/model"
	"github.com/tv-anagha/ecommerce-backend/cart-service/internal/repository"
	"github.com/tv-anagha/ecommerce-backend/cart-service/internal/service"
)

type CartHandler struct {
	svc *service.CartService
}

func NewCartHandler(svc *service.CartService) *CartHandler {
	return &CartHandler{svc: svc}
}

func (h *CartHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "cart-service"})
}

func parseUserID(c *gin.Context) (uint, bool) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return 0, false
	}
	return uint(userID), true
}

func parseProductID(c *gin.Context) (uint, bool) {
	productID, err := strconv.ParseUint(c.Param("productId"), 10, 64)
	if err != nil || productID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return 0, false
	}
	return uint(productID), true
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	items, err := h.svc.GetCart(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if items == nil {
		items = []model.CartItem{}
	}
	c.JSON(http.StatusOK, items)
}

type itemBody struct {
	ProductID uint `json:"productId"`
	Quantity  int  `json:"quantity" binding:"required"`
}

func (h *CartHandler) AddItem(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	var body itemBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.ProductID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "productId is required"})
		return
	}

	item, err := h.svc.AddItem(userID, body.ProductID, body.Quantity)
	if errors.Is(err, service.ErrInvalidQuantity) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if errors.Is(err, client.ErrProductNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

type quantityBody struct {
	Quantity int `json:"quantity" binding:"required"`
}

func (h *CartHandler) UpdateItem(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}
	productID, ok := parseProductID(c)
	if !ok {
		return
	}

	var body quantityBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.svc.UpdateItem(userID, productID, body.Quantity)
	if errors.Is(err, service.ErrInvalidQuantity) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if errors.Is(err, repository.ErrCartItemNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *CartHandler) RemoveItem(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}
	productID, ok := parseProductID(c)
	if !ok {
		return
	}

	if err := h.svc.RemoveItem(userID, productID); err != nil {
		if errors.Is(err, repository.ErrCartItemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
