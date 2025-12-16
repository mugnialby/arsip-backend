package request

type NewArchiveCharacteristicRequest struct {
	ArchiveCharacteristicName string `json:"archiveCharacteristicName" binding:"required"`
	CreatedBy                 string `json:"createdBy"`
}
