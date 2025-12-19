package request

type DeleteDepartmentRequest struct {
	ID          uint   `json:"id"`
	SubmittedBy string `json:"submittedBy"`
}
