package request

type NewArchiveRoleAccessRequest struct {
	ID           uint   `json:"id"`
	ArchiveID    uint   `json:"archiveId" binding:"required"`
	RoleID       uint   `json:"roleId" binding:"required"`
	DepartmentID uint   `json:"departmentId" binding:"required"`
	IsNew        bool   `json:"isNew"`
	IsDelete     bool   `json:"isDelete"`
	SubmittedBy  string `json:"submittedBy"`
}
