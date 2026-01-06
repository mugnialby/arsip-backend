package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mugnialby/arsip-backend/internal/model"
	request "github.com/mugnialby/arsip-backend/internal/model/dto/request/users"
	"github.com/mugnialby/arsip-backend/internal/service"
	"github.com/mugnialby/arsip-backend/pkg/logger"
	"github.com/mugnialby/arsip-backend/pkg/response"
	"go.uber.org/zap"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	users, err := h.service.GetAllUsers()
	if err != nil {
		logger.Log.Error("user.get_all.failed",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to get data")
		return
	}

	logger.Log.Info("user.get_all.success",
		zap.String("request_id", requestID.(string)),
		zap.Int("count", len(users)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	response.Success(c, users)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Log.Warn("user.get_by_id.invalid_id",
			zap.String("request_id", requestID.(string)),
			zap.String("param", c.Param("id")),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		logger.Log.Info("user.get_by_id.failed",
			zap.String("request_id", requestID.(string)),
			zap.Uint("user_id", uint(id)),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusNotFound, "Failed to get data")
		return
	}

	logger.Log.Info("user.get_by_id.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	response.Success(c, user)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	var newUserRequest request.NewUserRequest
	if err := c.ShouldBindJSON(&newUserRequest); err != nil {
		logger.Log.Warn("user.create.invalid_request",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
		)

		response.Error(c, http.StatusBadRequest, "JSON request is not valid")
		return
	}

	newUser := model.User{
		ID:           0,
		UserId:       newUserRequest.UserId,
		PasswordHash: newUserRequest.PasswordHash,
		FullName:     newUserRequest.FullName,
		DepartmentID: newUserRequest.DepartmentID,
		RoleID:       newUserRequest.RoleID,
		Status:       "Y",
		CreatedBy:    newUserRequest.SubmittedBy,
	}

	if err := h.service.CreateUser(&newUser); err != nil {
		logger.Log.Error("user.create.failed",
			zap.String("request_id", requestID.(string)),
			zap.Any("payload", newUser),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to create data")
		return
	}

	logger.Log.Info("user.create.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	c.Status(http.StatusCreated)
}

func (h *UserHandler) UpdateUserById(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	var updateUserRequest request.UpdateUserRequest
	if err := c.ShouldBindJSON(&updateUserRequest); err != nil {
		logger.Log.Warn("user.update.invalid_request",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
		return
	}

	user, err := h.service.GetUserByID(updateUserRequest.ID)
	if err != nil {
		logger.Log.Error("user.update.get_by_id.failed",
			zap.String("request_id", requestID.(string)),
			zap.Any("payload", updateUserRequest.ID),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusNotFound, "Failed to get data")
		return
	}

	timeNow := time.Now()
	user.UserId = updateUserRequest.UserId
	user.FullName = updateUserRequest.FullName
	user.DepartmentID = updateUserRequest.DepartmentID
	user.RoleID = updateUserRequest.RoleID
	user.ModifiedBy = &updateUserRequest.SubmittedBy
	user.ModifiedAt = &timeNow

	if err := h.service.UpdateUser(user); err != nil {
		logger.Log.Error("user.update.save.failed",
			zap.String("request_id", requestID.(string)),
			zap.Any("payload", updateUserRequest),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to update data")
		return
	}

	logger.Log.Info("user.update.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	response.Success(c, user)
}

func (h *UserHandler) DeleteUserById(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	var deleteUserRequest request.DeleteUserRequest
	if err := c.ShouldBindJSON(&deleteUserRequest); err != nil {
		logger.Log.Warn("user.delete.invalid_request",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
		)

		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
		return
	}

	if err := h.service.DeleteUser(&deleteUserRequest); err != nil {
		logger.Log.Error("user.delete.failed",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
			zap.Any("payload", deleteUserRequest),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to delete data")
		return
	}

	logger.Log.Info("user.delete.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	c.Status(http.StatusOK)
}
