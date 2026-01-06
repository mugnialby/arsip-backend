package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	request "github.com/mugnialby/arsip-backend/internal/model/dto/request/auth"
	"github.com/mugnialby/arsip-backend/internal/service"
	"github.com/mugnialby/arsip-backend/pkg/logger"
	"github.com/mugnialby/arsip-backend/pkg/response"
	"go.uber.org/zap"
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
	start := time.Now()
	requestID, _ := c.Get("request_id")

	var loginRequest request.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		logger.Log.Warn("auth.login.invalid_request",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
		)

		response.Error(c, http.StatusBadRequest, "JSON request is not valid")
		return
	}

	user, err := h.userService.CheckUserLoginRequest(&loginRequest)
	if err != nil {
		logger.Log.Error("auth.check_user_login.failed",
			zap.String("request_id", requestID.(string)),
			zap.Any("payload", loginRequest),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	logger.Log.Info("auth.login.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	response.Success(c, user)
}
