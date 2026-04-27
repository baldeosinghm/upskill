package courses

import (
	"context"
	"errors"
	"fmt"

	"github.com/baldeosinghm/upskill/internal/users"
	"github.com/jackc/pgx/v5"
)

type Service struct {
	repo         *Repository
	usersService *users.Service
}

func NewService(repo *Repository, usersService *users.Service) *Service {
	return &Service{
		repo:         repo,
		usersService: usersService,
	}
}

// Improve security subtlety w/ sentinel error
var (
	ErrOwnerNotFound   = errors.New("owner not found")
	ErrOwnerNotTeacher = errors.New("owner is not a teacher")
	ErrCourseNotFound  = errors.New("course not found")
)

func (s *Service) Create(ctx context.Context, name, ownerID string) (*Course, error) {
	// Check owner's role is a teacher
	user, err := s.usersService.GetByID(ctx, ownerID)
	if errors.Is(err, users.ErrUserNotFound) {
		return nil, ErrOwnerNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("create course: lookup owner: %w", err)
	}
	if user.Role != "teacher" {
		return nil, ErrOwnerNotTeacher
	}

	// Once validated, allow teach to create course
	course, err := s.repo.Create(ctx, name, ownerID)
	// If some DB issue came back, return error
	if err != nil {
		return nil, fmt.Errorf("create course: %w", err)
	}
	return course, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*Course, error) {
	course, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCourseNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get course by id: %w", err)
	}
	return course, nil
}

func (s *Service) List(ctx context.Context) ([]Course, error) {
	courses, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list courses: %w", err)
	}
	return courses, nil
}
