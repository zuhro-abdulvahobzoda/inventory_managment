package repository

import (
	"context"
	"inventory-api/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (int64, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
}

type ItemRepository interface {
	CreateItem(ctx context.Context, userID int64, item models.Item) (int64, error)
	GetItem(ctx context.Context, id int64) (*models.Item, error)
	ListItems(ctx context.Context) ([]models.Item, error)
	UpdateItem(ctx context.Context, userID int64, item models.Item) error
	DeleteItem(ctx context.Context, userID, id int64) error
	ItemHistory(ctx context.Context, itemID int64) ([]models.HistoryEntry, error)
}

type LocationRepository interface {
	CreateLocation(ctx context.Context, location models.Location) (int64, error)
	ListLocations(ctx context.Context) ([]models.Location, error)
}
