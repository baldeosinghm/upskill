package users

import (
	"context"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Create a user
func (service Service) CreateUser(ctx context.Context, username, email, passwordHash, role string) (*User, error) {
	user, err := service.repo.Create(ctx, username, email, passwordHash, role)
	if err != nil {
		return nil, err
	}
	return user, nil
}
