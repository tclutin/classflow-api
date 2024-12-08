package schedule

import (
	"github.com/tclutin/classflow-api/internal/domain/edu"
)

type DetailsScheduleDTO struct {
	Type        string
	SubjectName string
	Teacher     string
	Room        string
	IsEven      bool
	DayOfWeek   int
	StartTime   string
	EndTime     string
	Building    edu.Building
}

type FilterDTO struct {
	IsEven string
}
