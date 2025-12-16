package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model"
	request "github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model/dto/request/users"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/service"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/pkg/response"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, users)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var newUserRequest request.NewUserRequest
	if err := c.ShouldBindJSON(&newUserRequest); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
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
		CreatedBy:    newUserRequest.CreatedBy,
	}

	if err := h.service.CreateUser(&newUser); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusInternalServerError, "API Fail")
		return
	}

	c.Status(http.StatusCreated)
}

func (h *UserHandler) UpdateUserById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rack not found"})
		return
	}

	var updateUserRequest request.UpdateUserRequest
	if err := c.ShouldBind(&updateUserRequest); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
		return
	}

	user.UserId = updateUserRequest.UserId
	user.FullName = updateUserRequest.FullName
	user.DepartmentID = updateUserRequest.DepartmentID
	user.RoleID = updateUserRequest.RoleID
	user.ModifiedBy = &updateUserRequest.ModifiedBy

	if err := h.service.UpdateUser(user); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusInternalServerError, "API Fail")
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUserById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user.Status = "N"
	if err := h.service.UpdateUser(user); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusInternalServerError, "API Fail")
		return
	}

	c.JSON(http.StatusOK, user)
}
