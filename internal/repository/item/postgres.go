package item

import (
	"context"
	"errors"
	"inventory-api/internal/models"
	"inventory-api/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	foreignKeyViolationCode = "23503"
	itemColumns             = "id, name, quantity, location_id"
)

var _ repository.ItemRepository = (*itemRepository)(nil)

type itemRepository struct {
	pool *pgxpool.Pool
}

func NewItemRepository(pool *pgxpool.Pool) *itemRepository {
	return &itemRepository{pool: pool}
}

// querier is what both the pool and a transaction provide, so getOne works with either.
type querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func getOne(ctx context.Context, q querier, query string, args ...any) (*models.Item, error) {
	rows, err := q.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	m, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrItemNotFound
		}
		return nil, err
	}

	item := m.toModel()
	return &item, nil
}

func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == foreignKeyViolationCode
}
