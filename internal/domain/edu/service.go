package edu

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	domainErr "github.com/tclutin/classflow-api/internal/domain/errors"
)

type Repository interface {
	GetAllFaculty(ctx context.Context) ([]Faculty, error)
	GetAllProgramsByFacultyId(ctx context.Context, facultyID uint64) ([]Program, error)
	GetAllTypesOfSubject(ctx context.Context) ([]TypeOfSubject, error)
	GetAllBuildings(ctx context.Context) ([]Building, error)
	GetBuildingById(ctx context.Context, buildingID uint64) (Building, error)
	GetTypeOfSubjectById(ctx context.Context, typeOfSubjectId uint64) (TypeOfSubject, error)
	GetProgramById(ctx context.Context, programID uint64) (Program, error)
	GetFacultyById(ctx context.Context, facultyID uint64) (Faculty, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo}
}

func (s *Service) GetAllTypesOfSubject(ctx context.Context) ([]TypeOfSubject, error) {
	return s.repo.GetAllTypesOfSubject(ctx)
}

func (s *Service) GetAllFaculties(ctx context.Context) ([]Faculty, error) {
	return s.repo.GetAllFaculty(ctx)
}

func (s *Service) GetAllProgramsByFacultyId(ctx context.Context, facultyID uint64) ([]Program, error) {
	return s.repo.GetAllProgramsByFacultyId(ctx, facultyID)
}

func (s *Service) GetAllBuildings(ctx context.Context) ([]Building, error) {
	return s.repo.GetAllBuildings(ctx)
}

func (s *Service) GetBuildingById(ctx context.Context, buildingID uint64) (Building, error) {
	building, err := s.repo.GetBuildingById(ctx, buildingID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return building, domainErr.ErrBuildingNotFound
		}

		return building, fmt.Errorf("failed to get building: %w", err)
	}

	return building, nil
}

func (s *Service) GetTypeOfSubjectById(ctx context.Context, typeOfSubjectId uint64) (TypeOfSubject, error) {
	typeOfSubject, err := s.repo.GetTypeOfSubjectById(ctx, typeOfSubjectId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return TypeOfSubject{}, domainErr.ErrTypeOfSubjectNotFound
		}

		return TypeOfSubject{}, fmt.Errorf("failted to get type of subject: %w", err)
	}

	return typeOfSubject, nil
}

func (s *Service) GetProgramById(ctx context.Context, programID uint64) (Program, error) {
	program, err := s.repo.GetProgramById(ctx, programID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return program, domainErr.ErrProgramNotFound
		}

		return program, fmt.Errorf("failted to get program: %w", err)
	}

	return program, nil
}

func (s *Service) GetFacultyById(ctx context.Context, facultyID uint64) (Faculty, error) {
	faculty, err := s.repo.GetFacultyById(ctx, facultyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return faculty, domainErr.ErrFacultyNotFound
		}

		return faculty, fmt.Errorf("failted to get faculty: %w", err)
	}

	return faculty, nil
}
