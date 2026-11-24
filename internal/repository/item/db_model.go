package item

import (
	"encoding/json"
	"inventory-api/internal/models"
	"time"
)

type dbModel struct {
	ID         int64  `db:"id"`
	Name       string `db:"name"`
	Quantity   int64  `db:"quantity"`
	LocationID *int64 `db:"location_id"`
}

func (m dbModel) toModel() models.Item {
	return models.Item{
		ID:         m.ID,
		Name:       m.Name,
		Quantity:   m.Quantity,
		LocationID: m.LocationID,
	}
}

type historyDBModel struct {
	ID        int64     `db:"id"`
	ItemID    int64     `db:"item_id"`
	UserID    int64     `db:"user_id"`
	Action    string    `db:"action"`
	ChangedAt time.Time `db:"changed_at"`
	OldData   []byte    `db:"old_data"`
	NewData   []byte    `db:"new_data"`
}

func (m historyDBModel) toModel() models.HistoryEntry {
	return models.HistoryEntry{
		ID:        m.ID,
		ItemID:    m.ItemID,
		UserID:    m.UserID,
		Action:    m.Action,
		ChangedAt: m.ChangedAt,
		OldData:   m.OldData,
		NewData:   m.NewData,
	}
}

// snapshot is the item state stored in history.old_data / history.new_data.
type snapshot struct {
	Name       string `json:"name"`
	Quantity   int64  `json:"quantity"`
	LocationID *int64 `json:"location_id"`
}

// toJSON returns nil when there is no item (e.g. old_data on creation): it is stored as NULL in the database.
func toJSON(item *models.Item) ([]byte, error) {
	if item == nil {
		return nil, nil
	}

	return json.Marshal(snapshot{
		Name:       item.Name,
		Quantity:   item.Quantity,
		LocationID: item.LocationID,
	})
}
