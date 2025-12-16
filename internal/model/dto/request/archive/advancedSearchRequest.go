package request

import "time"

type AdvancedSearchRequest struct {
	ArchiveName   *string    `json:"archiveName"`
	ArchiveTypeID *uint      `json:"archiveTypeId"`
	ArchiveDate   *time.Time `json:"archiveDate"`
}
