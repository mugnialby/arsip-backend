package request

type NewRoleRequest struct {
	RoleName     string `json:"roleName" binding:"required"`
	DepartmentID uint   `json:"departmentId" binding:"required"`
	CreatedBy    string `json:"createdBy"`
}
