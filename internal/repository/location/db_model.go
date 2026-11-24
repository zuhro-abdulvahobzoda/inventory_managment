package location

import "inventory-api/internal/models"

type dbModel struct {
	ID        int64    `db:"id"`
	Name      string   `db:"name"`
	Latitude  *float64 `db:"latitude"`
	Longitude *float64 `db:"longitude"`
}

func (m dbModel) toModel() models.Location {
	return models.Location{
		ID:        m.ID,
		Name:      m.Name,
		Latitude:  m.Latitude,
		Longitude: m.Longitude,
	}
}
