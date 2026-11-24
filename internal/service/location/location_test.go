package location

import (
	"context"
	"errors"
	"inventory-api/internal/models"
	"inventory-api/internal/repository"
	"strings"
	"testing"
)

type fakeRepo struct {
	created []models.Location
}

var _ repository.LocationRepository = (*fakeRepo)(nil)

func (r *fakeRepo) CreateLocation(_ context.Context, location models.Location) (int64, error) {
	r.created = append(r.created, location)
	return int64(len(r.created)), nil
}

func (r *fakeRepo) ListLocations(context.Context) ([]models.Location, error) {
	return r.created, nil
}

func ptr(v float64) *float64 { return &v }

func TestCreateLocation(t *testing.T) {
	tests := []struct {
		name     string
		location models.Location
		wantErr  error
	}{
		{name: "success", location: models.Location{Name: "Warehouse A", Latitude: ptr(38.5598), Longitude: ptr(68.787)}},
		{name: "zero coordinates are valid", location: models.Location{Name: "Null Island", Latitude: ptr(0), Longitude: ptr(0)}},
		{name: "boundary values are valid", location: models.Location{Name: "Edge", Latitude: ptr(-90), Longitude: ptr(180)}},
		{name: "blank name", location: models.Location{Name: "  ", Latitude: ptr(1), Longitude: ptr(1)}, wantErr: models.ErrInvalidLocationName},
		{name: "name too long", location: models.Location{Name: strings.Repeat("a", 101), Latitude: ptr(1), Longitude: ptr(1)}, wantErr: models.ErrInvalidLocationName},
		{name: "missing latitude", location: models.Location{Name: "A", Longitude: ptr(1)}, wantErr: models.ErrInvalidCoordinates},
		{name: "missing longitude", location: models.Location{Name: "A", Latitude: ptr(1)}, wantErr: models.ErrInvalidCoordinates},
		{name: "latitude too big", location: models.Location{Name: "A", Latitude: ptr(90.1), Longitude: ptr(1)}, wantErr: models.ErrInvalidCoordinates},
		{name: "latitude too small", location: models.Location{Name: "A", Latitude: ptr(-90.1), Longitude: ptr(1)}, wantErr: models.ErrInvalidCoordinates},
		{name: "longitude too big", location: models.Location{Name: "A", Latitude: ptr(1), Longitude: ptr(180.1)}, wantErr: models.ErrInvalidCoordinates},
		{name: "longitude too small", location: models.Location{Name: "A", Latitude: ptr(1), Longitude: ptr(-180.1)}, wantErr: models.ErrInvalidCoordinates},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{}
			svc := NewLocationService(repo)

			_, err := svc.CreateLocation(context.Background(), tt.location)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if wantStored := tt.wantErr == nil; (len(repo.created) == 1) != wantStored {
				t.Fatalf("stored %d locations, want stored = %v", len(repo.created), wantStored)
			}
		})
	}
}
