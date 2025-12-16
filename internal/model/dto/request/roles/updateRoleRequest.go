package request

type UpdateRoleRequest struct {
	RoleName     string `json:"roleName" binding:"required"`
	DepartmentID uint   `json:"departmentID" binding:"required"`
	ModifiedBy   string `json:"modifiedBy"`
}
