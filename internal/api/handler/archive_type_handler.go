package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model"
	request "github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model/dto/request/archiveType"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/service"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/pkg/response"
)

type ArchiveTypeHandler struct {
	service *service.ArchiveTypeService
}

func NewArchiveTypeHandler(s *service.ArchiveTypeService) *ArchiveTypeHandler {
	return &ArchiveTypeHandler{service: s}
}

func (h *ArchiveTypeHandler) GetAllArchiveTypes(c *gin.Context) {
	archiveType, err := h.service.GetAllArchiveTypes()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, archiveType)
}

func (h *ArchiveTypeHandler) GetArchiveTypeByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	archiveType, err := h.service.GetArchiveTypeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "archive type not found"})
		return
	}
	c.JSON(http.StatusOK, archiveType)
}

func (h *ArchiveTypeHandler) CreateArchiveType(c *gin.Context) {
	var newArchiveTypeRequest request.NewArchiveTypeRequest
	if err := c.ShouldBindJSON(&newArchiveTypeRequest); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
		return
	}

	newArchiveType := model.ArchiveType{
		ID:              0,
		ArchiveTypeName: newArchiveTypeRequest.ArchiveTypeName,
		Status:          "Y",
		CreatedBy:       newArchiveTypeRequest.CreatedBy,
	}

	if err := h.service.CreateArchiveType(&newArchiveType); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusInternalServerError, "API Fail")
		return
	}

	c.Status(http.StatusCreated)
}

func (h *ArchiveTypeHandler) UpdateArchiveTypeById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	archiveType, err := h.service.GetArchiveTypeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "archiveType not found"})
		return
	}

	var updateArchiveTypeRequest request.UpdateArchiveTypeRequest
	if err := c.ShouldBind(&updateArchiveTypeRequest); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
		return
	}

	archiveType.ArchiveTypeName = updateArchiveTypeRequest.ArchiveTypeName
	if err := h.service.UpdateArchiveType(archiveType); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusInternalServerError, "API Fail")
		return
	}

	c.JSON(http.StatusOK, archiveType)
}

func (h *ArchiveTypeHandler) DeleteArchiveTypeById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	archiveType, err := h.service.GetArchiveTypeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "archiveType not found"})
		return
	}

	archiveType.Status = "N"
	if err := h.service.UpdateArchiveType(archiveType); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusInternalServerError, "API Fail")
		return
	}

	c.JSON(http.StatusOK, archiveType)
}
