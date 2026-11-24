package models

import "time"

const (
	ActionAdded   = "added"
	ActionChanged = "changed"
	ActionDeleted = "deleted"
)

type HistoryEntry struct {
	ID        int64
	ItemID    int64
	UserID    int64
	Action    string
	ChangedAt time.Time
	OldData   []byte
	NewData   []byte
}
