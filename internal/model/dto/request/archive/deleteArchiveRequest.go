package request

type DeleteArchiveRequest struct {
	ID     uint   `json:"id"`
	Status string `json:"status"`
}
