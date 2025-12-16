package request

type UpdateDepartmentRequest struct {
	DepartmentName string `json:"departmentName" binding:"required"`
}
