package request

type UpdateDepartmentRequest struct {
	ID             uint   `json:"id"`
	DepartmentName string `json:"departmentName" binding:"required"`
	SubmittedBy    string `json:"submittedBy"`
}
