package request

type UpdateArchiveCharacteristicRequest struct {
	ID                        uint   `json:"id"`
	ArchiveCharacteristicName string `json:"archiveCharacteristicName" binding:"required"`
	SubmittedBy               string `json:"submittedBy"`
}
