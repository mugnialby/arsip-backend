package request

import (
	request "github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model/dto/request/archiveAttachment"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/utils"
)

type NewArchiveRequest struct {
	ArchiveName             string                                `json:"archiveName" binding:"required"`
	ArchiveCharacteristicID uint                                  `json:"archiveCharacteristicId" binding:"required"`
	ArchiveTypeID           uint                                  `json:"archiveTypeId" binding:"required"`
	ArchiveDate             utils.DateOnly                        `json:"archiveDate"`
	ListArchiveAttachments  []request.NewArchiveAttachmentRequest `json:"listArchiveAttachments"`
	UserId                  string                                `json:"userId"`
}
