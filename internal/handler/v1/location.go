package v1

import (
	"inventory-api/internal/handler/response"
	"inventory-api/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createLocationRequest struct {
	Name string `json:"name" binding:"required"`
	// Pointers let us tell "not sent" apart from a valid 0 (equator / prime meridian).
	Latitude  *float64 `json:"latitude" binding:"required"`
	Longitude *float64 `json:"longitude" binding:"required"`
}

type createLocationResponse struct {
	ID int64 `json:"id"`
}

type locationResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	// Null only for locations created before coordinates were introduced.
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

func (h *Handler) createLocation(c *gin.Context) {
	var req createLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	id, err := h.srvs.LocationService.CreateLocation(c.Request.Context(), models.Location{
		Name:      req.Name,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	})
	if err != nil {
		writeError(c, "create location", err)
		return
	}

	c.JSON(http.StatusCreated, createLocationResponse{ID: id})
}

func (h *Handler) listLocations(c *gin.Context) {
	locations, err := h.srvs.LocationService.ListLocations(c.Request.Context())
	if err != nil {
		writeError(c, "list locations", err)
		return
	}

	resp := make([]locationResponse, 0, len(locations))
	for _, l := range locations {
		resp = append(resp, locationResponse{
			ID:        l.ID,
			Name:      l.Name,
			Latitude:  l.Latitude,
			Longitude: l.Longitude,
		})
	}

	c.JSON(http.StatusOK, resp)
}
