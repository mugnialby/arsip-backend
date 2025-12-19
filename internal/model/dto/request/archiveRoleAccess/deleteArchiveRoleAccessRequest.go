package request

type DeleteArchiveRoleAccessRequest struct {
	ID          uint   `json:"id"`
	SubmittedBy string `json:"submittedBy"`
}
