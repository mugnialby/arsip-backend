package request

type UpdateUserRequest struct {
	UserId       string `json:"userId" binding:"required"`
	FullName     string `json:"fullName" binding:"required"`
	DepartmentID uint   `json:"departmentId" binding:"required"`
	RoleID       uint   `json:"roleId" binding:"required"`
	ModifiedBy   string `json:"modifiedBy"`
}
