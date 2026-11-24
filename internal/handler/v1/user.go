package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type userResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (h *Handler) getMe(c *gin.Context) {
	user, err := h.srvs.UserService.GetByID(c.Request.Context(), currentUserID(c))
	if err != nil {
		writeError(c, "get me", err)
		return
	}

	c.JSON(http.StatusOK, userResponse{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	})
}
