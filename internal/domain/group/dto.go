package group

import "time"

type CreateGroupDTO struct {
	FacultyID uint64
	ProgramID uint64
	ShortName string
}

type DetailsGroupDTO struct {
	GroupID        uint64
	LeaderID       *uint64
	Faculty        string
	Program        string
	ShortName      string
	NumberOfPeople int
	ExistsSchedule bool
	CreatedAt      time.Time
}

type SummaryGroupDTO struct {
	GroupID        uint64
	Faculty        string
	Program        string
	ShortName      string
	NumberOfPeople int
	ExistsSchedule bool
}

type FilterDTO struct {
	Faculty string
	Program string
}
