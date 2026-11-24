package models

type Item struct {
	ID         int64
	Name       string
	Quantity   int64
	LocationID *int64
}
