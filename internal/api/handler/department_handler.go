package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mugnialby/arsip-backend/internal/model"
	request "github.com/mugnialby/arsip-backend/internal/model/dto/request/department"
	"github.com/mugnialby/arsip-backend/internal/service"
	"github.com/mugnialby/arsip-backend/pkg/logger"
	"github.com/mugnialby/arsip-backend/pkg/response"
	"go.uber.org/zap"
)

type DepartmentHandler struct {
	service *service.DepartmentService
}

func NewDepartmentHandler(s *service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: s}
}

func (h *DepartmentHandler) GetAllDepartments(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	departments, err := h.service.GetAllDepartments()
	if err != nil {
		logger.Log.Error("department.get_all.failed",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to get data")
		return
	}

	logger.Log.Info("department.get_all.success",
		zap.String("request_id", requestID.(string)),
		zap.Int("count", len(departments)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	response.Success(c, departments)
}

func (h *DepartmentHandler) GetDepartmentByID(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Log.Warn("department.get_by_id.invalid_id",
			zap.String("request_id", requestID.(string)),
			zap.String("param", c.Param("id")),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	department, err := h.service.GetDepartmentByID(uint(id))
	if err != nil {
		logger.Log.Info("department.get_by_id.failed",
			zap.String("request_id", requestID.(string)),
			zap.Uint("department_id", uint(id)),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusNotFound, "Failed to get data")
		return
	}

	logger.Log.Info("department.get_by_id.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	response.Success(c, department)
}

func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	var newDepartmentRequest request.NewDepartmentRequest
	if err := c.ShouldBindJSON(&newDepartmentRequest); err != nil {
		logger.Log.Warn("department.create.invalid_request",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
		)

		response.Error(c, http.StatusBadRequest, "JSON request is not valid")
		return
	}

	newDepartment := model.Department{
		ID:             0,
		DepartmentName: newDepartmentRequest.DepartmentName,
		Status:         "Y",
		CreatedBy:      newDepartmentRequest.SubmittedBy,
	}

	if err := h.service.CreateDepartment(&newDepartment); err != nil {
		logger.Log.Error("department.create.failed",
			zap.String("request_id", requestID.(string)),
			zap.Any("payload", newDepartment),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to create data")
		return
	}

	logger.Log.Info("department.create.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	c.Status(http.StatusCreated)
}

func (h *DepartmentHandler) UpdateDepartmentById(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	var updateDepartmentRequest request.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&updateDepartmentRequest); err != nil {
		logger.Log.Warn("department.update.invalid_request",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
		return
	}

	department, err := h.service.GetDepartmentByID(updateDepartmentRequest.ID)
	if err != nil {
		logger.Log.Error("department.update.get_by_id.failed",
			zap.String("request_id", requestID.(string)),
			zap.Any("payload", updateDepartmentRequest.ID),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusNotFound, "Failed to get data")
		return
	}

	timeNow := time.Now()
	department.DepartmentName = updateDepartmentRequest.DepartmentName
	department.ModifiedBy = &updateDepartmentRequest.SubmittedBy
	department.ModifiedAt = &timeNow

	if err := h.service.UpdateDepartment(department); err != nil {
		logger.Log.Error("department.update.save.failed",
			zap.String("request_id", requestID.(string)),
			zap.Any("payload", updateDepartmentRequest),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to update data")
		return
	}

	logger.Log.Info("department.update.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	response.Success(c, department)
}

func (h *DepartmentHandler) DeleteDepartmentById(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	var deleteDepartmentRequest request.DeleteDepartmentRequest
	if err := c.ShouldBindJSON(&deleteDepartmentRequest); err != nil {
		logger.Log.Warn("department.delete.invalid_request",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
		)

		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
		return
	}

	if err := h.service.DeleteDepartment(&deleteDepartmentRequest); err != nil {
		logger.Log.Error("department.delete.failed",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
			zap.Any("payload", deleteDepartmentRequest),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to delete data")
		return
	}

	logger.Log.Info("department.delete.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	c.Status(http.StatusOK)
}
