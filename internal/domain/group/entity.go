package group

import "time"

type Group struct {
	GroupID        uint64
	LeaderID       uint64
	FacultyID      uint64
	ProgramID      uint64
	ShortName      string
	Code           string
	NumberOfPeople int
	ExistsSchedule bool
	CreatedAt      time.Time
}
