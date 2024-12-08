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

type TypeOfSubjectResponse struct {
	TypeOfSubjectID uint64 `json:"type_of_subject_id"`
	Name            string `json:"name"`
}

type BuildingResponse struct {
	BuildingID uint64  `json:"building_id"`
	Name       string  `json:"name"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Address    string  `json:"address"`
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

func EntitiesToTypesOfSubjectResponse(entities []edu.TypeOfSubject) []TypeOfSubjectResponse {
	var typeOfSubjects []TypeOfSubjectResponse
	for _, entity := range entities {
		faculty := TypeOfSubjectResponse{
			TypeOfSubjectID: entity.TypeOfSubjectID,
			Name:            entity.Name,
		}

		typeOfSubjects = append(typeOfSubjects, faculty)
	}

	return typeOfSubjects
}

func EntitiesToBuildingsResponse(entities []edu.Building) []BuildingResponse {
	var buildings []BuildingResponse
	for _, entity := range entities {
		faculty := BuildingResponse{
			BuildingID: entity.BuildingID,
			Name:       entity.Name,
			Latitude:   entity.Latitude,
			Longitude:  entity.Longitude,
			Address:    entity.Address,
		}

		buildings = append(buildings, faculty)
	}

	return buildings
}
