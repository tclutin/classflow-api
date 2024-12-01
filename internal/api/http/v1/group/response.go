package group

import (
	"github.com/tclutin/classflow-api/internal/api/http/v1/edu"
	"github.com/tclutin/classflow-api/internal/domain/group"
	"github.com/tclutin/classflow-api/internal/domain/schedule"
	"time"
)

type SummaryGroupResponse struct {
	GroupID        uint64 `json:"group_id"`
	Faculty        string `json:"faculty"`
	Program        string `json:"program"`
	ShortName      string `json:"short_name"`
	NumberOfPeople int    `json:"number_of_people"`
	ExistsSchedule bool   `json:"exists_schedule"`
}

type DetailsGroupResponse struct {
	GroupID        uint64    `json:"group_id"`
	LeaderID       *uint64   `json:"leader_id"`
	Faculty        string    `json:"faculty"`
	Program        string    `json:"program"`
	ShortName      string    `json:"short_name"`
	NumberOfPeople int       `json:"number_of_people"`
	ExistsSchedule bool      `json:"exists_schedule"`
	CreatedAt      time.Time `json:"created_at"`
}

type DetailsScheduleResponse struct {
	Type        string               `json:"type"`
	SubjectName string               `json:"subject_name"`
	Teacher     string               `json:"teacher"`
	Room        string               `json:"room"`
	IsEven      bool                 `json:"is_even"`
	DayOfWeek   int                  `json:"day_of_week"`
	StartTime   string               `json:"start_time"`
	EndTime     string               `json:"end_time"`
	Building    edu.BuildingResponse `json:"building"`
}

func EntitiesToSummaryGroupsResponse(entities []group.SummaryGroupDTO) []SummaryGroupResponse {
	var summaryGroupsResponse []SummaryGroupResponse
	for _, entity := range entities {
		summaryGroupResponse := SummaryGroupResponse{
			GroupID:        entity.GroupID,
			Faculty:        entity.Faculty,
			Program:        entity.Program,
			ShortName:      entity.ShortName,
			NumberOfPeople: entity.NumberOfPeople,
			ExistsSchedule: entity.ExistsSchedule,
		}

		summaryGroupsResponse = append(summaryGroupsResponse, summaryGroupResponse)
	}

	return summaryGroupsResponse

}

func EntityToDetailsGroupResponse(entity group.DetailsGroupDTO) DetailsGroupResponse {
	return DetailsGroupResponse{
		GroupID:        entity.GroupID,
		LeaderID:       entity.LeaderID,
		Faculty:        entity.Faculty,
		Program:        entity.Program,
		ShortName:      entity.ShortName,
		NumberOfPeople: entity.NumberOfPeople,
		ExistsSchedule: entity.ExistsSchedule,
		CreatedAt:      entity.CreatedAt,
	}

}

func EntitiesToSchedulesResponse(entities []schedule.DetailsScheduleDTO) []DetailsScheduleResponse {
	var schedulesResponse []DetailsScheduleResponse

	for _, entity := range entities {
		scheduleResponse := DetailsScheduleResponse{
			Type:        entity.Type,
			SubjectName: entity.SubjectName,
			Teacher:     entity.Teacher,
			Room:        entity.Room,
			IsEven:      entity.IsEven,
			DayOfWeek:   entity.DayOfWeek,
			StartTime:   entity.StartTime,
			EndTime:     entity.EndTime,
			Building: edu.BuildingResponse{
				BuildingID: entity.Building.BuildingID,
				Name:       entity.Building.Name,
				Latitude:   entity.Building.Latitude,
				Longitude:  entity.Building.Longitude,
				Address:    entity.Building.Address,
			},
		}

		schedulesResponse = append(schedulesResponse, scheduleResponse)
	}

	return schedulesResponse
}
