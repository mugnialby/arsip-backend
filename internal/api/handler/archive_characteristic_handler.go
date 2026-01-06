package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mugnialby/arsip-backend/internal/model"
	request "github.com/mugnialby/arsip-backend/internal/model/dto/request/archiveCharacteristic"
	"github.com/mugnialby/arsip-backend/internal/service"
	"github.com/mugnialby/arsip-backend/pkg/logger"
	"github.com/mugnialby/arsip-backend/pkg/response"
	"go.uber.org/zap"
)

type ArchiveCharacteristicHandler struct {
	service *service.ArchiveCharacteristicService
}

func NewArchiveCharacteristicHandler(s *service.ArchiveCharacteristicService) *ArchiveCharacteristicHandler {
	return &ArchiveCharacteristicHandler{service: s}
}

func (h *ArchiveCharacteristicHandler) GetAllArchiveCharacteristics(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	archiveCharacteristics, err := h.service.GetAllArchiveCharacteristics()
	if err != nil {
		logger.Log.Error("archive_characteristic.get_all.failed",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to get data")
		return
	}

	logger.Log.Info("archive_characteristic.get_all.success",
		zap.String("request_id", requestID.(string)),
		zap.Int("count", len(archiveCharacteristics)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	response.Success(c, archiveCharacteristics)
}

func (h *ArchiveCharacteristicHandler) GetArchiveCharacteristicByID(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Log.Warn("archive_characteristic.get_by_id.invalid_id",
			zap.String("request_id", requestID.(string)),
			zap.String("param", c.Param("id")),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	archiveCharacteristic, err := h.service.GetArchiveCharacteristicByID(uint(id))
	if err != nil {
		logger.Log.Info("archive_characteristic.get_by_id.failed",
			zap.String("request_id", requestID.(string)),
			zap.Uint("archive_characteristic_id", uint(id)),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusNotFound, "Failed to get data")
		return
	}

	logger.Log.Info("archive_characteristic.get_by_id.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	response.Success(c, archiveCharacteristic)
}

func (h *ArchiveCharacteristicHandler) CreateArchiveCharacteristic(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	var newArchiveCharacteristicRequest request.NewArchiveCharacteristicRequest
	if err := c.ShouldBindJSON(&newArchiveCharacteristicRequest); err != nil {
		logger.Log.Warn("archive_characteristic.create.invalid_request",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
		)

		response.Error(c, http.StatusBadRequest, "JSON request is not valid")
		return
	}

	newArchiveCharacteristic := model.ArchiveCharacteristic{
		ID:                        0,
		ArchiveCharacteristicName: newArchiveCharacteristicRequest.ArchiveCharacteristicName,
		Status:                    "Y",
		CreatedBy:                 newArchiveCharacteristicRequest.SubmittedBy,
	}

	if err := h.service.CreateArchiveCharacteristic(&newArchiveCharacteristic); err != nil {
		logger.Log.Error("archive_characteristic.create.failed",
			zap.String("request_id", requestID.(string)),
			zap.Any("payload", newArchiveCharacteristic),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to create data")
		return
	}

	logger.Log.Info("archive_characteristic.create.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	c.Status(http.StatusCreated)
}

func (h *ArchiveCharacteristicHandler) UpdateArchiveCharacteristicById(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	var updateArchiveCharacteristicRequest request.UpdateArchiveCharacteristicRequest
	if err := c.ShouldBindJSON(&updateArchiveCharacteristicRequest); err != nil {
		logger.Log.Warn("archive_characteristic.update.invalid_request",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
		return
	}

	archiveCharacteristic, err := h.service.GetArchiveCharacteristicByID(updateArchiveCharacteristicRequest.ID)
	if err != nil {
		logger.Log.Error("archive_characteristic.update.get_by_id.failed",
			zap.String("request_id", requestID.(string)),
			zap.Any("payload", updateArchiveCharacteristicRequest.ID),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusNotFound, "Failed to get data")
		return
	}

	timeNow := time.Now()
	archiveCharacteristic.ArchiveCharacteristicName = updateArchiveCharacteristicRequest.ArchiveCharacteristicName
	archiveCharacteristic.ModifiedBy = &updateArchiveCharacteristicRequest.SubmittedBy
	archiveCharacteristic.ModifiedAt = &timeNow

	if err := h.service.UpdateArchiveCharacteristic(archiveCharacteristic); err != nil {
		logger.Log.Error("archive_characteristic.update.save.failed",
			zap.String("request_id", requestID.(string)),
			zap.Any("payload", updateArchiveCharacteristicRequest),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to update data")
		return
	}

	logger.Log.Info("archive_characteristic.update.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	response.Success(c, archiveCharacteristic)
}

func (h *ArchiveCharacteristicHandler) DeleteArchiveCharacteristicById(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	var deleteArchiveCharacteristicRequest request.DeleteArchiveCharacteristicRequest
	if err := c.ShouldBindJSON(&deleteArchiveCharacteristicRequest); err != nil {
		logger.Log.Warn("archive_characteristic.delete.invalid_request",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
		)

		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
		return
	}

	if err := h.service.DeleteArchiveCharacteristic(&deleteArchiveCharacteristicRequest); err != nil {
		logger.Log.Error("archive_characteristic.delete.failed",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
			zap.Any("payload", deleteArchiveCharacteristicRequest),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to delete data")
	}

	logger.Log.Info("archive_characteristic.delete.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	c.Status(http.StatusOK)
}
