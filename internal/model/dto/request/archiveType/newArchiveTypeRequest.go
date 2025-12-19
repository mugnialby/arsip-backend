package request

type NewArchiveTypeRequest struct {
	ArchiveTypeName string `json:"archiveTypeName" binding:"required"`
	SubmittedBy     string `json:"submittedBy"`
}
