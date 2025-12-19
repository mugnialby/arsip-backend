package request

type DeleteArchiveTypeRequest struct {
	ID          uint   `json:"id"`
	SubmittedBy string `json:"submittedBy"`
}
