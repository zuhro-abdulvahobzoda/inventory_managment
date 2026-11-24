package location

import (
	"context"
	"fmt"
	"inventory-api/internal/models"

	"github.com/jackc/pgx/v5"
)

func (r *locationRepository) CreateLocation(ctx context.Context, location models.Location) (int64, error) {
	var id int64
	err := r.pool.QueryRow(ctx,
		`INSERT INTO locations (name, latitude, longitude) VALUES ($1, $2, $3) RETURNING id`,
		location.Name, location.Latitude, location.Longitude,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create location: %w", err)
	}

	return id, nil
}

func (r *locationRepository) ListLocations(ctx context.Context) ([]models.Location, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, latitude, longitude FROM locations ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("list locations: %w", err)
	}

	dbLocations, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbModel])
	if err != nil {
		return nil, fmt.Errorf("list locations: %w", err)
	}

	locations := make([]models.Location, 0, len(dbLocations))
	for _, m := range dbLocations {
		locations = append(locations, m.toModel())
	}

	return locations, nil
}
