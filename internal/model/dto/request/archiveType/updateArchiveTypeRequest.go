package request

type UpdateArchiveTypeRequest struct {
	ArchiveTypeName string `json:"archiveTypeName" binding:"required"`
	ModifiedBy      string `json:"modifiedBy"`
}
