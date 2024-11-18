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
