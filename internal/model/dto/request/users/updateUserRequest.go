package request

type UpdateUserRequest struct {
	ID           uint   `json:"id"`
	UserId       string `json:"userId" binding:"required"`
	FullName     string `json:"fullName" binding:"required"`
	DepartmentID uint   `json:"departmentId" binding:"required"`
	RoleID       uint   `json:"roleId" binding:"required"`
	SubmittedBy  string `json:"submittedBy"`
}
