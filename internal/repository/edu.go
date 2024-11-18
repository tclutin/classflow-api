package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tclutin/classflow-api/internal/domain/edu"
)

type EduRepository struct {
	pool *pgxpool.Pool
}

func NewEduRepository(pool *pgxpool.Pool) *EduRepository {
	return &EduRepository{pool}
}

func (f *EduRepository) GetAllProgramsByFacultyId(ctx context.Context, facultyID uint64) ([]edu.Program, error) {
	sql := `SELECT * FROM public.programs WHERE faculty_id = $1`

	rows, err := f.pool.Query(ctx, sql, facultyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var programs []edu.Program

	for rows.Next() {
		var program edu.Program
		if err = rows.Scan(&program.ProgramID, &program.FacultyID, &program.Name); err != nil {
			return nil, err
		}

		programs = append(programs, program)
	}

	return programs, nil
}

func (f *EduRepository) GetAllFaculty(ctx context.Context) ([]edu.Faculty, error) {
	sql := `SELECT * FROM public.faculties`

	rows, err := f.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var faculties []edu.Faculty

	for rows.Next() {
		var faculty edu.Faculty
		if err = rows.Scan(&faculty.FacultyID, &faculty.Name); err != nil {
			return nil, err
		}

		faculties = append(faculties, faculty)
	}

	return faculties, nil
}

func (f *EduRepository) GetAllBuildings(ctx context.Context) ([]edu.Building, error) {
	sql := `SELECT * FROM public.buildings`

	rows, err := f.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buildings []edu.Building

	for rows.Next() {
		var building edu.Building
		if err = rows.Scan(&building.BuildingID, &building.Name, &building.Latitude, &building.Longitude, &building.Address); err != nil {
			return nil, err
		}

		buildings = append(buildings, building)
	}

	return buildings, nil
}

func (f *EduRepository) GetBuildingById(ctx context.Context, buildingID uint64) (edu.Building, error) {
	sql := `SELECT * FROM public.buildings WHERE buildings_id = $1`

	row := f.pool.QueryRow(ctx, sql, buildingID)

	var building edu.Building

	if err := row.Scan(&building.BuildingID, &building.Name, &building.Latitude, &building.Longitude, &building.Address); err != nil {
		return building, err
	}

	return building, nil
}

func (f *EduRepository) GetAllTypesOfSubject(ctx context.Context) ([]edu.TypeOfSubject, error) {
	sql := `SELECT * FROM public.type_of_subject`

	rows, err := f.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var typesOfSubject []edu.TypeOfSubject

	for rows.Next() {
		var typeOfSubject edu.TypeOfSubject
		if err = rows.Scan(&typeOfSubject.TypeOfSubjectID, &typeOfSubject.Name); err != nil {
			return nil, err
		}

		typesOfSubject = append(typesOfSubject, typeOfSubject)
	}

	return typesOfSubject, nil
}

func (f *EduRepository) GetTypeOfSubjectById(ctx context.Context, typeOfSubjectId uint64) (edu.TypeOfSubject, error) {
	sql := `SELECT * FROM public.type_of_subject WHERE type_of_subject_id = $1`

	row := f.pool.QueryRow(ctx, sql, typeOfSubjectId)

	var typeOfSubject edu.TypeOfSubject

	err := row.Scan(&typeOfSubject.TypeOfSubjectID, &typeOfSubject.Name)
	if err != nil {
		return edu.TypeOfSubject{}, err
	}

	return typeOfSubject, nil
}

func (f *EduRepository) GetFacultyById(ctx context.Context, facultyID uint64) (edu.Faculty, error) {
	sql := `SELECT * FROM public.faculties WHERE faculty_id = $1`

	row := f.pool.QueryRow(ctx, sql, facultyID)

	var faculty edu.Faculty

	err := row.Scan(&faculty.FacultyID, &faculty.Name)
	if err != nil {
		return edu.Faculty{}, err
	}

	return faculty, nil
}

func (f *EduRepository) GetProgramById(ctx context.Context, programID uint64) (edu.Program, error) {
	sql := `SELECT * FROM public.programs WHERE program_id = $1`

	row := f.pool.QueryRow(ctx, sql, programID)

	var program edu.Program

	err := row.Scan(&program.ProgramID, &program.FacultyID, &program.Name)
	if err != nil {
		return edu.Program{}, err
	}

	return program, nil
}
