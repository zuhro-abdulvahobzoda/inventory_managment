package v1

import (
	"inventory-api/internal/handler/response"
	"inventory-api/internal/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	bearerPrefix = "Bearer "
	userIDKey    = "userID"
	userRoleKey  = "userRole"
)

func (h *Handler) authMiddleware(c *gin.Context) {
	header := c.GetHeader("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(header, bearerPrefix))
	if !strings.HasPrefix(header, bearerPrefix) || token == "" {
		response.Error(c, http.StatusUnauthorized, "missing bearer token")
		return
	}

	user, err := h.srvs.AuthService.Authenticate(c.Request.Context(), token)
	if err != nil {
		writeError(c, "auth middleware", err)
		return
	}

	c.Set(userIDKey, user.ID)
	c.Set(userRoleKey, user.Role)
	c.Next()
}

func (h *Handler) adminMiddleware(c *gin.Context) {
	if c.GetString(userRoleKey) != models.RoleAdmin {
		writeError(c, "admin middleware", models.ErrForbidden)
		return
	}

	c.Next()
}

func currentUserID(c *gin.Context) int64 {
	return c.GetInt64(userIDKey)
}
