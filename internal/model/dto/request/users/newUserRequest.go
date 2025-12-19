package request

type NewUserRequest struct {
	UserId       string `json:"userId" binding:"required"`
	PasswordHash string `json:"passwordHash" binding:"required"`
	FullName     string `json:"fullName" binding:"required"`
	DepartmentID uint   `json:"departmentId" binding:"required"`
	RoleID       uint   `json:"roleId" binding:"required"`
	SubmittedBy  string `json:"submittedBy"`
}
