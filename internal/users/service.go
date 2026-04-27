package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Improve security subtlety w/ sentinel error
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserNotFound = errors.New("user not found")

// Create a user
func (s *Service) CreateUser(ctx context.Context, username, email, password, role string) (*User, error) {
	// Hash user's password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, fmt.Errorf("unable to hash password: %w", err)
	}

	// Pass hashed password to db
	user, err := s.repo.Create(ctx, username, email, string(passwordHash), role)

	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.Login(ctx, email)

	// User doesn't exist, return error
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrInvalidCredentials
	}
	// If some DB issue came back, return error
	if err != nil {
		return "", fmt.Errorf("server error: %w", err)
	}

	// Unhash password and compare to user-entered password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	// Incorrect password, return error
	if err != nil {
		return "", ErrInvalidCredentials
	}

	// Password matches to email, return user's ID
	return user.ID, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*User, error) {
	user, err := s.repo.GetByID(ctx, id)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}
