package edu

import "github.com/tclutin/classflow-api/internal/domain/edu"

type FacultyResponse struct {
	FacultyID uint64 `json:"faculty_id"`
	Name      string `json:"name"`
}

type ProgramResponse struct {
	ProgramID uint64 `json:"program_id"`
	FacultyID uint64 `json:"faculty_id"`
	Name      string `json:"name"`
}

func EntitiesToFacultiesResponse(entities []edu.Faculty) []FacultyResponse {
	var faculties []FacultyResponse

	for _, entity := range entities {
		faculty := FacultyResponse{
			FacultyID: entity.FacultyID,
			Name:      entity.Name,
		}
		faculties = append(faculties, faculty)
	}

	return faculties
}

func EntitiesToProgramsResponse(entities []edu.Program) []ProgramResponse {
	var faculties []ProgramResponse

	for _, entity := range entities {
		faculty := ProgramResponse{
			ProgramID: entity.ProgramID,
			FacultyID: entity.FacultyID,
			Name:      entity.Name,
		}
		faculties = append(faculties, faculty)
	}

	return faculties
}
