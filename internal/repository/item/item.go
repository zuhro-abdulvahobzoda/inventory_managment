package item

import (
	"context"
	"fmt"
	"inventory-api/internal/models"

	"github.com/jackc/pgx/v5"
)

// Creating, updating, deleting happens in transaction with history:
// or the item is saved with its history or nothing

func (r *itemRepository) CreateItem(ctx context.Context, userID int64, item models.Item) (int64, error) {
	var id int64

	err := pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			`INSERT INTO items (name, quantity, location_id) VALUES ($1, $2, $3) RETURNING id`,
			item.Name, item.Quantity, item.LocationID,
		).Scan(&id)
		if err != nil {
			return err
		}

		item.ID = id
		return addHistory(ctx, tx, userID, id, models.ActionAdded, nil, &item)
	})
	if err != nil {
		if isForeignKeyViolation(err) {
			return 0, models.ErrLocationNotFound
		}
		return 0, fmt.Errorf("create item: %w", err)
	}

	return id, nil
}

func (r *itemRepository) GetItem(ctx context.Context, id int64) (*models.Item, error) {
	item, err := getOne(ctx, r.pool,
		`SELECT `+itemColumns+` FROM items WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return nil, fmt.Errorf("get item: %w", err)
	}

	return item, nil
}

func (r *itemRepository) ListItems(ctx context.Context) ([]models.Item, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+itemColumns+` FROM items WHERE deleted_at IS NULL ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}

	dbItems, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbModel])
	if err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}

	items := make([]models.Item, 0, len(dbItems))
	for _, m := range dbItems {
		items = append(items, m.toModel())
	}

	return items, nil
}

func (r *itemRepository) UpdateItem(ctx context.Context, userID int64, item models.Item) error {
	err := pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
		// FOR UPDATE locks the row so that old_data in the history exactly matches what we are overwriting.
		old, err := getOne(ctx, tx,
			`SELECT `+itemColumns+` FROM items WHERE id = $1 AND deleted_at IS NULL FOR UPDATE`, item.ID)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx,
			`UPDATE items SET name = $1, quantity = $2, location_id = $3 WHERE id = $4`,
			item.Name, item.Quantity, item.LocationID, item.ID,
		)
		if err != nil {
			return err
		}

		return addHistory(ctx, tx, userID, item.ID, models.ActionChanged, old, &item)
	})
	if err != nil {
		if isForeignKeyViolation(err) {
			return models.ErrLocationNotFound
		}
		return fmt.Errorf("update item: %w", err)
	}

	return nil
}

// DeleteItem — "soft deletion", the item still exists wiith Deleted_at value
func (r *itemRepository) DeleteItem(ctx context.Context, userID, id int64) error {
	err := pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
		old, err := getOne(ctx, tx,
			`SELECT `+itemColumns+` FROM items WHERE id = $1 AND deleted_at IS NULL FOR UPDATE`, id)
		if err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, `UPDATE items SET deleted_at = now() WHERE id = $1`, id); err != nil {
			return err
		}

		return addHistory(ctx, tx, userID, id, models.ActionDeleted, old, nil)
	})
	if err != nil {
		return fmt.Errorf("delete item: %w", err)
	}

	return nil
}
