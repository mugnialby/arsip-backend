package request

type DeleteArchiveCharacteristicRequest struct {
	ID          uint   `json:"id"`
	SubmittedBy string `json:"submittedBy"`
}
