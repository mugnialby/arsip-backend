package request

type GetArchiveByDataRequest struct {
	DepartmentID uint `json:"departmentId" binding:"required"`
	RoleID       uint `json:"roleId" binding:"required"`
}
