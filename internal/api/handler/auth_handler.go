package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	request "github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model/dto/request/auth"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/service"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/pkg/response"
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
