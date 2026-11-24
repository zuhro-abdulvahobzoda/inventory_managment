package item

import (
	"context"
	"fmt"
	"inventory-api/internal/models"

	"github.com/jackc/pgx/v5"
)

func addHistory(ctx context.Context, tx pgx.Tx, userID, itemID int64, action string, oldItem, newItem *models.Item) error {
	oldData, err := toJSON(oldItem)
	if err != nil {
		return fmt.Errorf("marshal old data: %w", err)
	}

	newData, err := toJSON(newItem)
	if err != nil {
		return fmt.Errorf("marshal new data: %w", err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO history (item_id, user_id, action, old_data, new_data) VALUES ($1, $2, $3, $4, $5)`,
		itemID, userID, action, oldData, newData,
	)
	if err != nil {
		return fmt.Errorf("add history: %w", err)
	}

	return nil
}

// ItemHistory returns the history of an item, including a deleted one.
// Every item has at least an "added" record, so an empty history means the item does not exist.
func (r *itemRepository) ItemHistory(ctx context.Context, itemID int64) ([]models.HistoryEntry, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, item_id, user_id, action, changed_at, old_data, new_data
		 FROM history WHERE item_id = $1 ORDER BY id`, itemID)
	if err != nil {
		return nil, fmt.Errorf("item history: %w", err)
	}

	dbEntries, err := pgx.CollectRows(rows, pgx.RowToStructByName[historyDBModel])
	if err != nil {
		return nil, fmt.Errorf("item history: %w", err)
	}

	if len(dbEntries) == 0 {
		return nil, models.ErrItemNotFound
	}

	entries := make([]models.HistoryEntry, 0, len(dbEntries))
	for _, m := range dbEntries {
		entries = append(entries, m.toModel())
	}

	return entries, nil
}
