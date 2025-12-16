package request

type NewArchiveTypeRequest struct {
	ArchiveTypeName string `json:"archiveTypeName" binding:"required"`
	CreatedBy       string `json:"createdBy"`
}
