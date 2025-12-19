package request

type UpdateRoleRequest struct {
	ID           uint   `json:"id"`
	RoleName     string `json:"roleName" binding:"required"`
	DepartmentID uint   `json:"departmentID" binding:"required"`
	SubmittedBy  string `json:"submittedBy"`
}
