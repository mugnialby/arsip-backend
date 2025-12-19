package request

type DeleteRoleRequest struct {
	ID          uint   `json:"id"`
	SubmittedBy string `json:"submittedBy"`
}
