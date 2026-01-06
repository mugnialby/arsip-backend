package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mugnialby/arsip-backend/internal/model"
	request "github.com/mugnialby/arsip-backend/internal/model/dto/request/archiveType"
	"github.com/mugnialby/arsip-backend/internal/service"
	"github.com/mugnialby/arsip-backend/pkg/logger"
	"github.com/mugnialby/arsip-backend/pkg/response"
	"go.uber.org/zap"
)

type ArchiveTypeHandler struct {
	service *service.ArchiveTypeService
}

func NewArchiveTypeHandler(s *service.ArchiveTypeService) *ArchiveTypeHandler {
	return &ArchiveTypeHandler{service: s}
}

func (h *ArchiveTypeHandler) GetAllArchiveTypes(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	archiveTypes, err := h.service.GetAllArchiveTypes()
	if err != nil {
		logger.Log.Error("archive_type.get_all.failed",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to get data")
		return
	}

	logger.Log.Info("archive_type.get_all.success",
		zap.String("request_id", requestID.(string)),
		zap.Int("count", len(archiveTypes)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	response.Success(c, archiveTypes)
}

func (h *ArchiveTypeHandler) GetArchiveTypeByID(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Log.Warn("archive_type.get_by_id.invalid_id",
			zap.String("request_id", requestID.(string)),
			zap.String("param", c.Param("id")),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	archiveType, err := h.service.GetArchiveTypeByID(uint(id))
	if err != nil {
		logger.Log.Info("archive_type.get_by_id.failed",
			zap.String("request_id", requestID.(string)),
			zap.Uint("archive_type_id", uint(id)),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusNotFound, "Failed to get data")
		return
	}

	logger.Log.Info("archive_type.get_by_id.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	response.Success(c, archiveType)
}

func (h *ArchiveTypeHandler) CreateArchiveType(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	var newArchiveTypeRequest request.NewArchiveTypeRequest
	if err := c.ShouldBindJSON(&newArchiveTypeRequest); err != nil {
		logger.Log.Warn("archive_type.create.invalid_request",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
		)

		response.Error(c, http.StatusBadRequest, "JSON request is not valid")
		return
	}

	newArchiveType := model.ArchiveType{
		ID:              0,
		ArchiveTypeName: newArchiveTypeRequest.ArchiveTypeName,
		Status:          "Y",
		CreatedBy:       newArchiveTypeRequest.SubmittedBy,
	}

	if err := h.service.CreateArchiveType(&newArchiveType); err != nil {
		logger.Log.Error("archive_type.create.failed",
			zap.String("request_id", requestID.(string)),
			zap.Any("payload", newArchiveType),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to create data")
		return
	}

	logger.Log.Info("archive_type.create.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	c.Status(http.StatusCreated)
}

func (h *ArchiveTypeHandler) UpdateArchiveTypeById(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	var updateArchiveTypeRequest request.UpdateArchiveTypeRequest
	if err := c.ShouldBindJSON(&updateArchiveTypeRequest); err != nil {
		logger.Log.Warn("archive_type.update.invalid_request",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
		return
	}

	archiveType, err := h.service.GetArchiveTypeByID(updateArchiveTypeRequest.ID)
	if err != nil {
		logger.Log.Error("archive_type.update.get_by_id.failed",
			zap.String("request_id", requestID.(string)),
			zap.Any("payload", updateArchiveTypeRequest.ID),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusNotFound, "Failed to get data")
		return
	}

	timeNow := time.Now()
	archiveType.ArchiveTypeName = updateArchiveTypeRequest.ArchiveTypeName
	archiveType.ModifiedBy = &updateArchiveTypeRequest.SubmittedBy
	archiveType.ModifiedAt = &timeNow

	if err := h.service.UpdateArchiveType(archiveType); err != nil {
		logger.Log.Error("archive_type.update.save.failed",
			zap.String("request_id", requestID.(string)),
			zap.Any("payload", updateArchiveTypeRequest),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to update data")
		return
	}

	logger.Log.Info("archive_type.update.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	response.Success(c, archiveType)
}

func (h *ArchiveTypeHandler) DeleteArchiveTypeById(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	var deleteArchiveTypeRequest request.DeleteArchiveTypeRequest
	if err := c.ShouldBindJSON(&deleteArchiveTypeRequest); err != nil {
		logger.Log.Warn("archive_type.delete.invalid_request",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
		)

		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
		return
	}

	if err := h.service.DeleteArchiveType(&deleteArchiveTypeRequest); err != nil {
		logger.Log.Error("archive_type.delete.failed",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
			zap.Any("payload", deleteArchiveTypeRequest),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to delete data")
		return
	}

	logger.Log.Info("archive_type.delete.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	c.Status(http.StatusOK)
}
