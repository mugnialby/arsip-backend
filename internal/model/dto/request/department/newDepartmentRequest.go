package request

type NewDepartmentRequest struct {
	DepartmentName string `json:"departmentName" binding:"required"`
	CreatedBy      string `json:"createdBy"`
}
