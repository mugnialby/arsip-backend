package request

type UpdateArchiveCharacteristicRequest struct {
	ArchiveCharacteristicName string `json:"archiveCharacteristicName" binding:"required"`
	ModifiedBy                string `json:"modifiedBy"`
}
