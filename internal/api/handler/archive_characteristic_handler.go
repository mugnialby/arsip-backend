package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model"
	request "github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model/dto/request/archiveCharacteristic"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/service"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/pkg/response"
)

type ArchiveCharacteristicHandler struct {
	service *service.ArchiveCharacteristicService
}

func NewArchiveCharacteristicHandler(s *service.ArchiveCharacteristicService) *ArchiveCharacteristicHandler {
	return &ArchiveCharacteristicHandler{service: s}
}

func (h *ArchiveCharacteristicHandler) GetAllArchiveCharacteristics(c *gin.Context) {
	archiveCharacteristic, err := h.service.GetAllArchiveCharacteristics()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, archiveCharacteristic)
}

func (h *ArchiveCharacteristicHandler) GetArchiveCharacteristicByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	archiveCharacteristic, err := h.service.GetArchiveCharacteristicByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "archive characteristic not found"})
		return
	}
	c.JSON(http.StatusOK, archiveCharacteristic)
}

func (h *ArchiveCharacteristicHandler) CreateArchiveCharacteristic(c *gin.Context) {
	var newArchiveCharacteristicRequest request.NewArchiveCharacteristicRequest
	if err := c.ShouldBindJSON(&newArchiveCharacteristicRequest); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
		return
	}

	newArchiveCharacteristic := model.ArchiveCharacteristic{
		ID:                        0,
		ArchiveCharacteristicName: newArchiveCharacteristicRequest.ArchiveCharacteristicName,
		Status:                    "Y",
		CreatedBy:                 newArchiveCharacteristicRequest.CreatedBy}

	if err := h.service.CreateArchiveCharacteristic(&newArchiveCharacteristic); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusInternalServerError, "API Fail")
		return
	}

	c.Status(http.StatusCreated)
}

func (h *ArchiveCharacteristicHandler) UpdateArchiveCharacteristicById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	archiveCharacteristic, err := h.service.GetArchiveCharacteristicByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "archive characteristic not found"})
		return
	}

	var updateArchiveCharacteristicRequest request.UpdateArchiveCharacteristicRequest
	if err := c.ShouldBind(&updateArchiveCharacteristicRequest); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
		return
	}

	archiveCharacteristic.ArchiveCharacteristicName = updateArchiveCharacteristicRequest.ArchiveCharacteristicName
	// role.ModifiedBy = &updateArchiveCharacteristicRequest.UserID
	if err := h.service.UpdateArchiveCharacteristic(archiveCharacteristic); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusInternalServerError, "API Fail")
		return
	}

	c.JSON(http.StatusOK, archiveCharacteristic)
}

func (h *ArchiveCharacteristicHandler) DeleteArchiveCharacteristicById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	archiveCharacteristic, err := h.service.GetArchiveCharacteristicByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "archiveCharacteristic not found"})
		return
	}

	archiveCharacteristic.Status = "N"
	if err := h.service.UpdateArchiveCharacteristic(archiveCharacteristic); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusInternalServerError, "API Fail")
		return
	}

	c.JSON(http.StatusOK, archiveCharacteristic)
}
