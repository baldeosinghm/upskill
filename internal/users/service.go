package users

import (
	"context"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Create a user
func (service Service) CreateUser(ctx context.Context, username, email, password, role string) (*User, error) {
	// Hash user's password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Pass hashed password to db
	user, err := service.repo.Create(ctx, username, email, string(passwordHash), role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service Service) Login(ctx context.Context, email, password string) (bool, error) {
	verified := false
	passwordHash, err := service.repo.Login(ctx, email)
	if err != nil {
		return verified, nil
	}
	// Unhash password and compare to user-entered password
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))

	// Return bool (will relay message of valid/invalid password/credentials)
	if err != nil {
		return verified, err
	}
	return verified, nil
}

func (service Service) LocateEmail(ctx context.Context, email string) (bool, error) {
	exists, err := service.repo.FindByEmail(ctx, email)
	if err != nil {
		return exists, err
	}
	log.Println("user exists")
	return exists, nil
}
