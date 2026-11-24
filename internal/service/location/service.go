package location

import (
	"inventory-api/internal/repository"
	"inventory-api/internal/service"
)

var _ service.LocationService = (*locationService)(nil)

type locationService struct {
	locationRepo repository.LocationRepository
}

func NewLocationService(locationRepo repository.LocationRepository) service.LocationService {
	return &locationService{locationRepo: locationRepo}
}
