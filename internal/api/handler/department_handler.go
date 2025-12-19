package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model"
	request "github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model/dto/request/department"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/service"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/pkg/response"
)

type DepartmentHandler struct {
	service *service.DepartmentService
}

func NewDepartmentHandler(s *service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: s}
}

func (h *DepartmentHandler) GetAllDepartments(c *gin.Context) {
	departments, err := h.service.GetAllDepartments()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, departments)
}

func (h *DepartmentHandler) GetDepartmentByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	department, err := h.service.GetDepartmentByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "department not found"})
		return
	}
	c.JSON(http.StatusOK, department)
}

func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	var newDepartmentRequest request.NewDepartmentRequest
	if err := c.ShouldBindJSON(&newDepartmentRequest); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
		return
	}

	newDepartment := model.Department{
		ID:             0,
		DepartmentName: newDepartmentRequest.DepartmentName,
		Status:         "Y",
		CreatedBy:      newDepartmentRequest.SubmittedBy,
	}

	if err := h.service.CreateDepartment(&newDepartment); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusInternalServerError, "API Fail")
		return
	}

	c.Status(http.StatusCreated)
}

func (h *DepartmentHandler) UpdateDepartmentById(c *gin.Context) {
	var updateDepartmentRequest request.UpdateDepartmentRequest
	if err := c.ShouldBind(&updateDepartmentRequest); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
		return
	}

	department, err := h.service.GetDepartmentByID(updateDepartmentRequest.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "data not found"})
		return
	}

	timeNow := time.Now()

	department.DepartmentName = updateDepartmentRequest.DepartmentName
	department.ModifiedBy = &updateDepartmentRequest.SubmittedBy
	department.ModifiedAt = &timeNow

	if err := h.service.UpdateDepartment(department); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusInternalServerError, "API Fail")
		return
	}

	c.JSON(http.StatusOK, department)
}

func (h *DepartmentHandler) DeleteDepartmentById(c *gin.Context) {
	var deleteDepartmentRequest request.DeleteDepartmentRequest
	if err := c.ShouldBindJSON(&deleteDepartmentRequest); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
		return
	}

	if err := h.service.DeleteDepartment(&deleteDepartmentRequest); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusInternalServerError, "API Fail")
		return
	}

	c.Status(http.StatusOK)
}
