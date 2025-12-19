package request

type NewDepartmentRequest struct {
	DepartmentName string `json:"departmentName" binding:"required"`
	SubmittedBy    string `json:"submittedBy"`
}
