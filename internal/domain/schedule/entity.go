package schedule

import "time"

type Schedule struct {
	ScheduleID      uint64
	groupID         uint64
	buildingsID     uint64
	TypeOfSubjectID int64
	SubjectName     string
	Room            string
	IsEven          bool
	DayOfWeek       int
	StartTime       string
	EndTime         string
	CreatedAt       time.Time
}
