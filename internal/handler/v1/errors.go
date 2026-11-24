package v1

import (
	"errors"
	"inventory-api/internal/handler/response"
	"inventory-api/internal/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func writeError(c *gin.Context, op string, err error) {
	switch {
	case errors.Is(err, models.ErrInvalidInput):
		response.Error(c, http.StatusBadRequest, models.ErrInvalidInput.Error())

	case errors.Is(err, models.ErrInvalidUsername):
		response.Error(c, http.StatusBadRequest, models.ErrInvalidUsername.Error())

	case errors.Is(err, models.ErrPasswordTooShort):
		response.Error(c, http.StatusBadRequest, models.ErrPasswordTooShort.Error())

	case errors.Is(err, models.ErrPasswordTooLong):
		response.Error(c, http.StatusBadRequest, models.ErrPasswordTooLong.Error())

	case errors.Is(err, models.ErrUsernameAlreadyExists):
		response.Error(c, http.StatusConflict, models.ErrUsernameAlreadyExists.Error())

	case errors.Is(err, models.ErrInvalidRole):
		response.Error(c, http.StatusBadRequest, models.ErrInvalidRole.Error())

	case errors.Is(err, models.ErrAdminAlreadyExists):
		response.Error(c, http.StatusConflict, models.ErrAdminAlreadyExists.Error())

	case errors.Is(err, models.ErrInvalidCredentials):
		response.Error(c, http.StatusUnauthorized, models.ErrInvalidCredentials.Error())

	case errors.Is(err, models.ErrInvalidToken):
		response.Error(c, http.StatusUnauthorized, models.ErrInvalidToken.Error())

	case errors.Is(err, models.ErrForbidden):
		response.Error(c, http.StatusForbidden, models.ErrForbidden.Error())

	case errors.Is(err, models.ErrUserNotFound):
		response.Error(c, http.StatusNotFound, models.ErrUserNotFound.Error())

	case errors.Is(err, models.ErrItemNotFound):
		response.Error(c, http.StatusNotFound, models.ErrItemNotFound.Error())

	case errors.Is(err, models.ErrInvalidItemName):
		response.Error(c, http.StatusBadRequest, models.ErrInvalidItemName.Error())

	case errors.Is(err, models.ErrInvalidQuantity):
		response.Error(c, http.StatusBadRequest, models.ErrInvalidQuantity.Error())

	case errors.Is(err, models.ErrLocationNotFound):
		response.Error(c, http.StatusBadRequest, models.ErrLocationNotFound.Error())

	case errors.Is(err, models.ErrInvalidLocationName):
		response.Error(c, http.StatusBadRequest, models.ErrInvalidLocationName.Error())

	case errors.Is(err, models.ErrInvalidCoordinates):
		response.Error(c, http.StatusBadRequest, models.ErrInvalidCoordinates.Error())
	default:

		log.Printf("handler %s: %v", op, err)
		response.Error(c, http.StatusInternalServerError, "internal server error")
	}
}
