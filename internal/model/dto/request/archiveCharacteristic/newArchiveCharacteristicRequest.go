package request

type NewArchiveCharacteristicRequest struct {
	ArchiveCharacteristicName string `json:"archiveCharacteristicName" binding:"required"`
	SubmittedBy               string `json:"submittedBy"`
}
