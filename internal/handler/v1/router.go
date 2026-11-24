package v1

import (
	"inventory-api/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	srvs *service.Services
}

func NewHandler(srvs *service.Services) *Handler {
	return &Handler{
		srvs: srvs,
	}
}

func (h *Handler) InitRoutes(rg *gin.RouterGroup) {

	rg.POST("/sign-up", h.signUp)
	rg.POST("/login", h.login)

	private := rg.Group("", h.authMiddleware)
	{
		private.GET("/users/me", h.getMe)

		private.GET("/items", h.listItems)
		private.GET("/items/:id", h.getItem)
		private.GET("/items/:id/history", h.itemHistory)

		private.GET("/locations", h.listLocations)

		admin := private.Group("", h.adminMiddleware)
		{
			admin.POST("/items", h.createItem)
			admin.PUT("/items/:id", h.updateItem)
			admin.DELETE("/items/:id", h.deleteItem)

			admin.POST("/locations", h.createLocation)
		}
	}
}
