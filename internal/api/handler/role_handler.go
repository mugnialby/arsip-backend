package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mugnialby/arsip-backend/internal/model"
	request "github.com/mugnialby/arsip-backend/internal/model/dto/request/roles"
	"github.com/mugnialby/arsip-backend/internal/service"
	"github.com/mugnialby/arsip-backend/pkg/response"
)

type RoleHandler struct {
	service *service.RoleService
}

func NewRoleHandler(s *service.RoleService) *RoleHandler {
	return &RoleHandler{service: s}
}

func (h *RoleHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.service.GetAllRoles()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, roles)
}

func (h *RoleHandler) GetRoleByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	role, err := h.service.GetRoleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}
	c.JSON(http.StatusOK, role)
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	var newRoleRequest request.NewRoleRequest
	if err := c.ShouldBindJSON(&newRoleRequest); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
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
		// tambah logger di sini
		response.Error(c, http.StatusInternalServerError, "API Fail")
		return
	}

	c.Status(http.StatusCreated)
}

func (h *RoleHandler) UpdateRoleById(c *gin.Context) {
	var updateRoleRequest request.UpdateRoleRequest
	if err := c.ShouldBind(&updateRoleRequest); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
		return
	}

	role, err := h.service.GetRoleByID(updateRoleRequest.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "data not found"})
		return
	}

	timeNow := time.Now()

	role.RoleName = updateRoleRequest.RoleName
	role.DepartmentID = updateRoleRequest.DepartmentID
	role.ModifiedBy = &updateRoleRequest.SubmittedBy
	role.ModifiedAt = &timeNow
	if err := h.service.UpdateRole(role); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusInternalServerError, "API Fail")
		return
	}

	c.JSON(http.StatusOK, role)
}

func (h *RoleHandler) DeleteRoleById(c *gin.Context) {
	var deleteRoleRequest request.DeleteRoleRequest
	if err := c.ShouldBindJSON(&deleteRoleRequest); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
		return
	}

	if err := h.service.DeleteRole(&deleteRoleRequest); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusInternalServerError, "API Fail")
		return
	}

	c.Status(http.StatusOK)
}

func (h *RoleHandler) GetRoleByDepartmentID(c *gin.Context) {
	departmentId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bad Request"})
		return
	}

	roles, err := h.service.GetRoleByDepartmentID(uint(departmentId))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Roles not found"})
		return
	}

	c.JSON(http.StatusOK, roles)
}
