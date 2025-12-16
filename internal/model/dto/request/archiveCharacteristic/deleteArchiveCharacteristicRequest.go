package request

type DeleteArchiveCharacteristicRequest struct {
	ID     uint   `json:"id"`
	Status string `json:"status"`
}
