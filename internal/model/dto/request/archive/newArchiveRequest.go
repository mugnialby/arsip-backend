package request

import (
	attachmentRequest "github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model/dto/request/archiveAttachment"
	roleAccessRequest "github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model/dto/request/archiveRoleAccess"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/utils"
)

type NewArchiveRequest struct {
	ArchiveDate             utils.DateOnly                                  `json:"archiveDate"`
	ArchiveNumber           string                                          `json:"archiveNumber"`
	ArchiveName             string                                          `json:"archiveName" binding:"required"`
	ArchiveCharacteristicID uint                                            `json:"archiveCharacteristicId" binding:"required"`
	ArchiveTypeID           uint                                            `json:"archiveTypeId" binding:"required"`
	DepartmentID            uint                                            `json:"departmentId" binding:"required"`
	ListArchiveAttachments  []attachmentRequest.NewArchiveAttachmentRequest `json:"listArchiveAttachments"`
	RoleAccess              []roleAccessRequest.NewArchiveRoleAccessRequest `json:"roleAccess"`
	SubmittedBy             string                                          `json:"submittedBy"`
}
