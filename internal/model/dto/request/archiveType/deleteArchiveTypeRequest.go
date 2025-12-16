package request

type DeleteArchiveTypeRequest struct {
	ID     uint   `json:"id"`
	Status string `json:"status"`
}
