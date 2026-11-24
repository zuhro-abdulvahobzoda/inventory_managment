package location

import (
	"context"
	"fmt"
	"inventory-api/internal/models"
	"math"
	"strings"
	"unicode/utf8"
)

const (
	maxNameLength = 100
	maxLatitude   = 90
	maxLongitude  = 180
)

func (s *locationService) CreateLocation(ctx context.Context, location models.Location) (int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	location.Name = strings.TrimSpace(location.Name)
	if location.Name == "" || utf8.RuneCountInString(location.Name) > maxNameLength {
		return 0, models.ErrInvalidLocationName
	}

	if !validCoordinate(location.Latitude, maxLatitude) || !validCoordinate(location.Longitude, maxLongitude) {
		return 0, models.ErrInvalidCoordinates
	}

	id, err := s.locationRepo.CreateLocation(ctx, location)
	if err != nil {
		return 0, fmt.Errorf("failed to create location: %w", err)
	}

	return id, nil
}

func (s *locationService) ListLocations(ctx context.Context) ([]models.Location, error) {
	locations, err := s.locationRepo.ListLocations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list locations: %w", err)
	}

	return locations, nil
}

// validCoordinate reports whether v is set and lies within [-limit, limit].
func validCoordinate(v *float64, limit float64) bool {
	return v != nil && !math.IsNaN(*v) && math.Abs(*v) <= limit
}
