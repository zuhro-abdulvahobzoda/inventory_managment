package models

type Location struct {
	ID   int64
	Name string
	// Latitude and Longitude are in decimal degrees (WGS 84).
	Latitude  *float64
	Longitude *float64
}
