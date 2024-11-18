package group

type CreateGroupRequest struct {
	FacultyID uint64 `json:"faculty_id" binding:"required"`
	ProgramID uint64 `json:"program_id" binding:"required"`
	ShortName string `json:"short_name" binding:"required,min=4,max=12"`
}

type JoinToGroupRequest struct {
	Code string `json:"code" binding:"required,max=10"`
}

type UploadScheduleRequest struct{}
