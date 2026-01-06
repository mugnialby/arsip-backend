package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mugnialby/arsip-backend/internal/model"
	request "github.com/mugnialby/arsip-backend/internal/model/dto/request/roles"
	"github.com/mugnialby/arsip-backend/internal/service"
	"github.com/mugnialby/arsip-backend/pkg/logger"
	"github.com/mugnialby/arsip-backend/pkg/response"
	"go.uber.org/zap"
)

type RoleHandler struct {
	service *service.RoleService
}

func NewRoleHandler(s *service.RoleService) *RoleHandler {
	return &RoleHandler{service: s}
}

func (h *RoleHandler) GetAllRoles(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	roles, err := h.service.GetAllRoles()
	if err != nil {
		logger.Log.Error("role.get_all.failed",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to get data")
		return
	}

	logger.Log.Info("role.get_all.success",
		zap.String("request_id", requestID.(string)),
		zap.Int("count", len(roles)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	response.Success(c, roles)
}

func (h *RoleHandler) GetRoleByID(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Log.Warn("role.get_by_id.invalid_id",
			zap.String("request_id", requestID.(string)),
			zap.String("param", c.Param("id")),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	role, err := h.service.GetRoleByID(uint(id))
	if err != nil {
		logger.Log.Info("role.get_by_id.failed",
			zap.String("request_id", requestID.(string)),
			zap.Uint("role_id", uint(id)),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusNotFound, "Failed to get data")
		return
	}

	logger.Log.Info("role.get_by_id.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	response.Success(c, role)
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	var newRoleRequest request.NewRoleRequest
	if err := c.ShouldBindJSON(&newRoleRequest); err != nil {
		logger.Log.Warn("role.create.invalid_request",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
		)

		response.Error(c, http.StatusBadRequest, "JSON request is not valid")
		return
	}

	newRole := model.Role{
		ID:           0,
		RoleName:     newRoleRequest.RoleName,
		DepartmentID: newRoleRequest.DepartmentID,
		Status:       "Y",
		CreatedBy:    newRoleRequest.CreatedBy,
	}

	if err := h.service.CreateRole(&newRole); err != nil {
		logger.Log.Error("role.create.failed",
			zap.String("request_id", requestID.(string)),
			zap.Any("payload", newRole),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to create data")
		return
	}

	logger.Log.Info("role.create.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	c.Status(http.StatusCreated)
}

func (h *RoleHandler) UpdateRoleById(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	var updateRoleRequest request.UpdateRoleRequest
	if err := c.ShouldBindJSON(&updateRoleRequest); err != nil {
		logger.Log.Warn("role.update.invalid_request",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
		return
	}

	role, err := h.service.GetRoleByID(updateRoleRequest.ID)
	if err != nil {
		logger.Log.Error("role.update.get_by_id.failed",
			zap.String("request_id", requestID.(string)),
			zap.Any("payload", updateRoleRequest.ID),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusNotFound, "Failed to get data")
		return
	}

	timeNow := time.Now()
	role.RoleName = updateRoleRequest.RoleName
	role.DepartmentID = updateRoleRequest.DepartmentID
	role.ModifiedBy = &updateRoleRequest.SubmittedBy
	role.ModifiedAt = &timeNow

	if err := h.service.UpdateRole(role); err != nil {
		logger.Log.Error("role.update.save.failed",
			zap.String("request_id", requestID.(string)),
			zap.Any("payload", updateRoleRequest),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to update data")
		return
	}

	logger.Log.Info("role.update.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	response.Success(c, role)
}

func (h *RoleHandler) DeleteRoleById(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	var deleteRoleRequest request.DeleteRoleRequest
	if err := c.ShouldBindJSON(&deleteRoleRequest); err != nil {
		logger.Log.Warn("role.delete.invalid_request",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
		)

		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
		return
	}

	if err := h.service.DeleteRole(&deleteRoleRequest); err != nil {
		logger.Log.Error("role.delete.failed",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
			zap.Any("payload", deleteRoleRequest),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to delete data")
		return
	}

	logger.Log.Info("role.delete.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	c.Status(http.StatusOK)
}

func (h *RoleHandler) GetRoleByDepartmentID(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	departmentId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Log.Warn("role.get_role_by_department_id.invalid_id",
			zap.String("request_id", requestID.(string)),
			zap.String("param", c.Param("id")),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	roles, err := h.service.GetRoleByDepartmentID(uint(departmentId))
	if err != nil {
		logger.Log.Info("role.get_role_by_department_id.failed",
			zap.String("request_id", requestID.(string)),
			zap.Uint("department_id", uint(departmentId)),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusNotFound, "Failed to get data")
		return
	}

	logger.Log.Info("role.get_role_by_department_id.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	response.Success(c, roles)
}
