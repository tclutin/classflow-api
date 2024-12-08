package schedule

import "time"

type Schedule struct {
	ScheduleID      uint64
	GroupID         uint64
	BuildingsID     uint64
	TypeOfSubjectID uint64
	SubjectName     string
	Teacher         string
	Room            string
	IsEven          bool
	DayOfWeek       int
	StartTime       string
	EndTime         string
	CreatedAt       time.Time
}
