package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model"
	request "github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model/dto/request/roles"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/service"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/pkg/response"
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
	id, _ := strconv.Atoi(c.Param("id"))
	role, err := h.service.GetRoleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}

	var updateRoleRequest request.UpdateRoleRequest
	if err := c.ShouldBind(&updateRoleRequest); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
		return
	}

	role.RoleName = updateRoleRequest.RoleName
	role.DepartmentID = updateRoleRequest.DepartmentID
	// role.ModifiedBy = &updateRoleRequest.UserID
	if err := h.service.UpdateRole(role); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusInternalServerError, "API Fail")
		return
	}

	c.JSON(http.StatusOK, role)
}

func (h *RoleHandler) DeleteRoleById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	role, err := h.service.GetRoleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}

	role.Status = "N"
	if err := h.service.UpdateRole(role); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusInternalServerError, "API Fail")
		return
	}

	c.JSON(http.StatusOK, role)
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
