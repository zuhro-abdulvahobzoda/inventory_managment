package v1

import (
	"inventory-api/internal/handler/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type signUpRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	// Role is optional: "admin" or "employee" (default). Only one admin can exist.
	Role string `json:"role"`
}

type signUpResponse struct {
	ID int64 `json:"id"`
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) signUp(c *gin.Context) {
	var req signUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	userID, err := h.srvs.AuthService.SignUp(c.Request.Context(), req.Username, req.Password, req.Role)
	if err != nil {
		writeError(c, "sign up", err)
		return
	}

	c.JSON(http.StatusCreated, signUpResponse{ID: userID})
}

func (h *Handler) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	tokens, err := h.srvs.AuthService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		writeError(c, "login", err)
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}
