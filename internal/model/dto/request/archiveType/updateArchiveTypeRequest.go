package request

type UpdateArchiveTypeRequest struct {
	ID              uint   `json:"id"`
	ArchiveTypeName string `json:"archiveTypeName" binding:"required"`
	SubmittedBy     string `json:"submittedBy"`
}
