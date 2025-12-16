package request

type DeleteDepartmentRequest struct {
	ID     uint   `json:"id"`
	Status string `json:"status"`
}
