package group

import (
	"github.com/tclutin/classflow-api/internal/domain/group"
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
	LeaderID       uint64    `json:"leader_id"`
	Faculty        string    `json:"faculty"`
	Program        string    `json:"program"`
	ShortName      string    `json:"short_name"`
	Code           string    `json:"code"`
	NumberOfPeople int       `json:"number_of_people"`
	ExistsSchedule bool      `json:"exists_schedule"`
	CreatedAt      time.Time `json:"created_at"`
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

func EntitiesToDetailsGroupsResponse(entities []group.DetailsGroupDTO) []DetailsGroupResponse {
	var detailsGroupsResponse []DetailsGroupResponse
	for _, entity := range entities {
		detailsGroupResponse := DetailsGroupResponse{
			GroupID:        entity.GroupID,
			LeaderID:       entity.LeaderID,
			Faculty:        entity.Faculty,
			Program:        entity.Program,
			ShortName:      entity.ShortName,
			Code:           entity.Code,
			NumberOfPeople: entity.NumberOfPeople,
			ExistsSchedule: entity.ExistsSchedule,
			CreatedAt:      entity.CreatedAt,
		}

		detailsGroupsResponse = append(detailsGroupsResponse, detailsGroupResponse)
	}

	return detailsGroupsResponse

}
