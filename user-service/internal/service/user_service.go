// Package service holds checkout auth rules (hash password, issue JWT).
package service

import (
	"errors"

	"github.com/tv-anagha/ecommerce-backend/user-service/internal/auth"
	"github.com/tv-anagha/ecommerce-backend/user-service/internal/model"
	"github.com/tv-anagha/ecommerce-backend/user-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Register stores a new shopper (email must be unique).
func (s *UserService) Register(email, password string) (*model.User, error) {
	if _, err := s.repo.FindByEmail(email); err == nil {
		return nil, repository.ErrEmailAlreadyUsed
	} else if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{Email: email, PasswordHash: string(hash)}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// Login checks the password and returns a JWT for checkout.
func (s *UserService) Login(email, password string) (token string, user *model.User, err error) {
	user, err = s.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", nil, repository.ErrUserNotFound
		}
		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, repository.ErrUserNotFound // same message as unknown email
	}

	token, err = auth.Sign(user.ID, user.Email)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

// Me loads the user behind a valid checkout token.
func (s *UserService) Me(userID uint) (*model.User, error) {
	return s.repo.FindByID(userID)
}
