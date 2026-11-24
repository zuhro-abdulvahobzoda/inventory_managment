package item

import (
	"inventory-api/internal/repository"
	"inventory-api/internal/service"
)

var _ service.ItemService = (*itemService)(nil)

type itemService struct {
	itemRepo repository.ItemRepository
}

func NewItemService(itemRepo repository.ItemRepository) service.ItemService {
	return &itemService{itemRepo: itemRepo}
}
