package request

type NewArchiveAttachmentRequest struct {
	ID         uint   `json:"id"`
	FileBase64 string `json:"fileBase64" binding:"required"`
	IsNew      bool   `json:"isNew"`
	IsDelete   bool   `json:"isDelete"`
}
