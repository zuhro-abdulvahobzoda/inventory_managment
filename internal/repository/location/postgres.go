package location

import (
	"inventory-api/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ repository.LocationRepository = (*locationRepository)(nil)

type locationRepository struct {
	pool *pgxpool.Pool
}

func NewLocationRepository(pool *pgxpool.Pool) *locationRepository {
	return &locationRepository{pool: pool}
}
