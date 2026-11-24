package service

import (
	"context"
	"inventory-api/internal/models"
)

type AuthService interface {
	SignUp(ctx context.Context, username, password, role string) (int64, error)
	Login(ctx context.Context, username, password string) (models.TokenPair, error)
	// Authenticate validates an access token and returns the user.
	Authenticate(ctx context.Context, token string) (*models.User, error)
}

type UserService interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
}

type ItemService interface {
	CreateItem(ctx context.Context, userID int64, item models.Item) (int64, error)
	GetItem(ctx context.Context, id int64) (*models.Item, error)
	ListItems(ctx context.Context) ([]models.Item, error)
	UpdateItem(ctx context.Context, userID int64, item models.Item) error
	DeleteItem(ctx context.Context, userID, id int64) error
	ItemHistory(ctx context.Context, itemID int64) ([]models.HistoryEntry, error)
}

type LocationService interface {
	CreateLocation(ctx context.Context, location models.Location) (int64, error)
	ListLocations(ctx context.Context) ([]models.Location, error)
}

type Services struct {
	AuthService     AuthService
	UserService     UserService
	ItemService     ItemService
	LocationService LocationService
}
