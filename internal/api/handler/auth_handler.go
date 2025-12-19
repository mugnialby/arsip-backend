package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	request "github.com/mugnialby/arsip-backend/internal/model/dto/request/auth"
	"github.com/mugnialby/arsip-backend/internal/service"
	"github.com/mugnialby/arsip-backend/pkg/response"
)

type AuthHandler struct {
	userService *service.UserService
}

func NewAuthHandler(
	userService *service.UserService,
) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginRequest request.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		response.Error(c, http.StatusBadRequest, "JSON is not valid")
		return
	}

	user, err := h.userService.CheckUserLoginRequest(&loginRequest)
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	c.JSON(http.StatusOK, user)
}
