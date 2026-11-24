package v1

import (
	"encoding/json"
	"inventory-api/internal/handler/response"
	"inventory-api/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type itemRequest struct {
	Name       string `json:"name" binding:"required"`
	Quantity   *int64 `json:"quantity" binding:"required"`
	LocationID *int64 `json:"location_id"`
}

type createItemResponse struct {
	ID int64 `json:"id"`
}

type itemResponse struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Quantity   int64  `json:"quantity"`
	LocationID *int64 `json:"location_id"`
}

type historyResponse struct {
	ID        int64           `json:"id"`
	ItemID    int64           `json:"item_id"`
	UserID    int64           `json:"user_id"`
	Action    string          `json:"action"`
	ChangedAt time.Time       `json:"changed_at"`
	OldData   json.RawMessage `json:"old_data"`
	NewData   json.RawMessage `json:"new_data"`
}

func toItemResponse(item models.Item) itemResponse {
	return itemResponse{
		ID:         item.ID,
		Name:       item.Name,
		Quantity:   item.Quantity,
		LocationID: item.LocationID,
	}
}

func itemID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, http.StatusBadRequest, "invalid item id")
		return 0, false
	}

	return id, true
}

func (h *Handler) createItem(c *gin.Context) {
	var req itemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	id, err := h.srvs.ItemService.CreateItem(c.Request.Context(), currentUserID(c), models.Item{
		Name:       req.Name,
		Quantity:   *req.Quantity,
		LocationID: req.LocationID,
	})
	if err != nil {
		writeError(c, "create item", err)
		return
	}

	c.JSON(http.StatusCreated, createItemResponse{ID: id})
}

func (h *Handler) listItems(c *gin.Context) {
	items, err := h.srvs.ItemService.ListItems(c.Request.Context())
	if err != nil {
		writeError(c, "list items", err)
		return
	}

	resp := make([]itemResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, toItemResponse(item))
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) getItem(c *gin.Context) {
	id, ok := itemID(c)
	if !ok {
		return
	}

	item, err := h.srvs.ItemService.GetItem(c.Request.Context(), id)
	if err != nil {
		writeError(c, "get item", err)
		return
	}

	c.JSON(http.StatusOK, toItemResponse(*item))
}

func (h *Handler) updateItem(c *gin.Context) {
	id, ok := itemID(c)
	if !ok {
		return
	}

	var req itemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.srvs.ItemService.UpdateItem(c.Request.Context(), currentUserID(c), models.Item{
		ID:         id,
		Name:       req.Name,
		Quantity:   *req.Quantity,
		LocationID: req.LocationID,
	})
	if err != nil {
		writeError(c, "update item", err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) deleteItem(c *gin.Context) {
	id, ok := itemID(c)
	if !ok {
		return
	}

	if err := h.srvs.ItemService.DeleteItem(c.Request.Context(), currentUserID(c), id); err != nil {
		writeError(c, "delete item", err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) itemHistory(c *gin.Context) {
	id, ok := itemID(c)
	if !ok {
		return
	}

	entries, err := h.srvs.ItemService.ItemHistory(c.Request.Context(), id)
	if err != nil {
		writeError(c, "item history", err)
		return
	}

	resp := make([]historyResponse, 0, len(entries))
	for _, e := range entries {
		resp = append(resp, historyResponse{
			ID:        e.ID,
			ItemID:    e.ItemID,
			UserID:    e.UserID,
			Action:    e.Action,
			ChangedAt: e.ChangedAt,
			OldData:   e.OldData,
			NewData:   e.NewData,
		})
	}

	c.JSON(http.StatusOK, resp)
}
