// Package handler turns HTTP requests into service calls and JSON responses.
package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tv-anagha/ecommerce-backend/user-service/internal/auth"
	"github.com/tv-anagha/ecommerce-backend/user-service/internal/repository"
	"github.com/tv-anagha/ecommerce-backend/user-service/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "user-service"})
}

type registerLoginBody struct {
	Email    string `json:"email"`
	Username string `json:"username"` // some UIs send username instead of email
	Password string `json:"password" binding:"required,min=8"`
}

func (b *registerLoginBody) identity() (string, error) {
	email := strings.TrimSpace(b.Email)
	if email == "" {
		email = strings.TrimSpace(b.Username)
	}
	if email == "" {
		return "", errors.New("email is required")
	}
	return email, nil
}

// Register creates an account before checkout.
func (h *UserHandler) Register(c *gin.Context) {
	var body registerLoginBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	email, err := body.identity()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.svc.Register(email, body.Password)
	if errors.Is(err, repository.ErrEmailAlreadyUsed) {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": user.ID, "email": user.Email})
}

// Login returns a JWT the frontend stores and sends when starting checkout.
func (h *UserHandler) Login(c *gin.Context) {
	var body registerLoginBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	email, err := body.identity()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.svc.Login(email, body.Password)
	if errors.Is(err, repository.ErrUserNotFound) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":  token,
		"userId": user.ID,
		"email":  user.Email,
	})
}

// Me proves the token is still valid and returns who is checking out.
func (h *UserHandler) Me(c *gin.Context) {
	userID, err := userIDFromAuthHeader(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
		return
	}

	user, err := h.svc.Me(userID)
	if errors.Is(err, repository.ErrUserNotFound) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": user.ID, "email": user.Email})
}

// userIDFromAuthHeader reads "Bearer <jwt>" and parses the user id claim.
func userIDFromAuthHeader(header string) (uint, error) {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return 0, errors.New("bad header")
	}
	id, _, err := auth.Parse(strings.TrimPrefix(header, prefix))
	return id, err
}
