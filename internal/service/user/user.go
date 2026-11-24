package user

import (
	"context"
	"fmt"
	"inventory-api/internal/models"
	"inventory-api/internal/repository"
	"inventory-api/internal/service"
)

var _ service.UserService = (*userService)(nil)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) service.UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) GetByID(ctx context.Context, id int64) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return user, nil
}
