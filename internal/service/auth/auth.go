package auth

import (
	"context"
	"errors"
	"fmt"
	"inventory-api/internal/models"
	"inventory-api/internal/repository"
	"inventory-api/internal/service"
	"inventory-api/internal/service/token"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

const (
	maxUsernameLength = 100
	minPasswordLength = 8
	maxPasswordBytes  = 72
)

var _ service.AuthService = (*authService)(nil)

type authService struct {
	users  repository.UserRepository
	tokens token.Manager
}

func NewAuthService(users repository.UserRepository, tokens token.Manager) service.AuthService {
	return &authService{
		users:  users,
		tokens: tokens,
	}
}

func (s *authService) SignUp(ctx context.Context, username, password, role string) (int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	username = strings.ToLower(strings.TrimSpace(username))

	if username == "" || password == "" {
		return 0, models.ErrInvalidInput
	}

	if utf8.RuneCountInString(username) > maxUsernameLength {
		return 0, models.ErrInvalidUsername
	}

	if utf8.RuneCountInString(password) < minPasswordLength {
		return 0, models.ErrPasswordTooShort
	}

	if len(password) > maxPasswordBytes {
		return 0, models.ErrPasswordTooLong
	}

	role = strings.ToLower(strings.TrimSpace(role))
	if role == "" {
		role = models.RoleEmployee
	}
	if !models.IsValidRole(role) {
		return 0, models.ErrInvalidRole
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user := &models.User{
		Username:     username,
		PasswordHash: string(passwordHash),
		Role:         role,
	}

	userID, err := s.users.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, models.ErrUsernameAlreadyExists) || errors.Is(err, models.ErrAdminAlreadyExists) {
			return 0, err
		}
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return userID, nil
}

func (s *authService) Login(ctx context.Context, username, password string) (models.TokenPair, error) {
	if err := ctx.Err(); err != nil {
		return models.TokenPair{}, err
	}

	username = strings.ToLower(strings.TrimSpace(username))
	if username == "" || password == "" {
		return models.TokenPair{}, models.ErrInvalidCredentials
	}

	user, err := s.users.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return models.TokenPair{}, models.ErrInvalidCredentials
		}

		return models.TokenPair{}, fmt.Errorf("failed to get user by username: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return models.TokenPair{}, models.ErrInvalidCredentials
	}

	accessToken, err := s.tokens.GenerateAccessToken(ctx, user.ID)
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.tokens.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) Authenticate(ctx context.Context, tokenValue string) (*models.User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if strings.TrimSpace(tokenValue) == "" {
		return nil, models.ErrInvalidToken
	}

	userID, err := s.tokens.Validate(ctx, tokenValue)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, models.ErrInvalidToken
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}
