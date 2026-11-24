package item

import (
	"context"
	"fmt"
	"inventory-api/internal/models"
	"strings"
	"unicode/utf8"
)

const maxNameLength = 200

func validate(item *models.Item) error {
	item.Name = strings.TrimSpace(item.Name)

	if item.Name == "" || utf8.RuneCountInString(item.Name) > maxNameLength {
		return models.ErrInvalidItemName
	}
	if item.Quantity < 0 {
		return models.ErrInvalidQuantity
	}

	return nil
}

func (s *itemService) CreateItem(ctx context.Context, userID int64, item models.Item) (int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	if err := validate(&item); err != nil {
		return 0, err
	}

	id, err := s.itemRepo.CreateItem(ctx, userID, item)
	if err != nil {
		return 0, fmt.Errorf("failed to create item: %w", err)
	}

	return id, nil
}

func (s *itemService) GetItem(ctx context.Context, id int64) (*models.Item, error) {
	item, err := s.itemRepo.GetItem(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	return item, nil
}

func (s *itemService) ListItems(ctx context.Context) ([]models.Item, error) {
	items, err := s.itemRepo.ListItems(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list items: %w", err)
	}

	return items, nil
}

func (s *itemService) UpdateItem(ctx context.Context, userID int64, item models.Item) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := validate(&item); err != nil {
		return err
	}

	if err := s.itemRepo.UpdateItem(ctx, userID, item); err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}

func (s *itemService) DeleteItem(ctx context.Context, userID, id int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := s.itemRepo.DeleteItem(ctx, userID, id); err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}

func (s *itemService) ItemHistory(ctx context.Context, itemID int64) ([]models.HistoryEntry, error) {
	entries, err := s.itemRepo.ItemHistory(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get item history: %w", err)
	}

	return entries, nil
}
