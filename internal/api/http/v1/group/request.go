package group

import (
	"errors"
	"github.com/tclutin/classflow-api/internal/domain/schedule"
	"time"
)

type CreateGroupRequest struct {
	FacultyID uint64 `json:"faculty_id" binding:"required,gte=1"`
	ProgramID uint64 `json:"program_id" binding:"required,gte=1"`
	ShortName string `json:"short_name" binding:"required,min=4,max=12"`
}

type SubjectRequest struct {
	Name       string `json:"name" binding:"required"`
	Room       string `json:"room" binding:"required"`
	Teacher    string `json:"teacher" binding:"required"`
	TypeID     uint64 `json:"type_id" binding:"required,gte=1"`
	BuildingID uint64 `json:"building_id" binding:"required,gte=1"`
	StartTime  string `json:"start_time" binding:"required"`
	EndTime    string `json:"end_time" binding:"required"`
}

type DaysRequest struct {
	DayNumber int              `json:"day_number" binding:"required"`
	Subjects  []SubjectRequest `json:"subjects" binding:"required"`
}

type WeekRequest struct {
	IsEven bool          `json:"is_even" binding:"required"`
	Days   []DaysRequest `json:"days" binding:"required"`
}

type UploadScheduleRequest struct {
	Weeks []WeekRequest `json:"weeks" binding:"required"`
}

// TODO: need to add validate of numbers of days
func (u UploadScheduleRequest) Validate() error {
	if len(u.Weeks) != 1 && len(u.Weeks) != 2 {
		return errors.New("UploadScheduleRequest wrong length of weeks ")
	}

	if len(u.Weeks) == 1 {
		if len(u.Weeks[0].Days) > 7 || len(u.Weeks[0].Days) < 1 {
			return errors.New("UploadScheduleRequest wrong length of days ")
		}
	}

	if len(u.Weeks) == 2 {
		if len(u.Weeks[0].Days) > 7 || len(u.Weeks[1].Days) > 7 {
			return errors.New("UploadScheduleRequest wrong length of days ")
		}

		if len(u.Weeks[1].Days) < 1 || len(u.Weeks[1].Days) < 1 {
			return errors.New("UploadScheduleRequest wrong length of days ")
		}

		if u.Weeks[0].IsEven == u.Weeks[1].IsEven {
			return errors.New("UploadScheduleRequest wrong even of weeks")
		}

	}

	return nil
}

func (u UploadScheduleRequest) TransformToEntities(groupID uint64) []schedule.Schedule {
	var schedules []schedule.Schedule

	for _, week := range u.Weeks {
		for _, day := range week.Days {
			for _, subject := range day.Subjects {
				schedule := schedule.Schedule{
					GroupID:         groupID,
					BuildingsID:     subject.BuildingID,
					TypeOfSubjectID: subject.TypeID,
					SubjectName:     subject.Name,
					Teacher:         subject.Teacher,
					Room:            subject.Room,
					IsEven:          week.IsEven,
					DayOfWeek:       day.DayNumber,
					StartTime:       subject.StartTime,
					EndTime:         subject.EndTime,
					CreatedAt:       time.Now(),
				}
				schedules = append(schedules, schedule)
			}
		}
	}
	return schedules
}
