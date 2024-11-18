package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tclutin/classflow-api/internal/domain/group"
)

type GroupRepository struct {
	pool *pgxpool.Pool
}

func NewGroupRepository(pool *pgxpool.Pool) *GroupRepository {
	return &GroupRepository{pool}
}

func (g *GroupRepository) Create(ctx context.Context, group group.Group) (uint64, error) {
	sql := `
	INSERT INTO public.groups
    (leader_id, faculty_id, program_id, code, short_name, exists_schedule, number_of_people, created_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING group_id;`

	row := g.pool.QueryRow(
		ctx,
		sql,
		group.LeaderID,
		group.FacultyID,
		group.ProgramID,
		group.Code,
		group.ShortName,
		group.ExistsSchedule,
		group.NumberOfPeople,
		group.CreatedAt)

	var groupId uint64

	if err := row.Scan(&groupId); err != nil {
		return 0, err
	}

	return groupId, nil
}

func (g *GroupRepository) Update(ctx context.Context, group group.Group) error {
	sql := `
		UPDATE
			public.groups
		SET
			leader_id = $1,
			faculty_id = $2,
			program_id = $3,
			short_name = $4,
			code = $5,
			number_of_people = $6,
			exists_schedule = $7,
			created_at = $8
		WHERE group_id = $9
		`

	_, err := g.pool.Exec(
		ctx,
		sql,
		group.LeaderID,
		group.FacultyID,
		group.ProgramID,
		group.ShortName,
		group.Code,
		group.NumberOfPeople,
		group.ExistsSchedule,
		group.CreatedAt,
		group.GroupID)

	return err
}

func (g *GroupRepository) GetSummaryGroups(ctx context.Context) ([]group.SummaryGroupDTO, error) {
	sql := `
		SELECT
		    g.group_id,
			f.faculty_name,
			p.program_name,
			g.short_name,
			g.number_of_people,
			g.exists_schedule
		FROM
			public.groups AS g
		INNER JOIN
			public.programs AS p ON g.program_id = p.program_id
		INNER JOIN
			public.faculties AS f ON g.faculty_id = f.faculty_id
	`

	rows, err := g.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []group.SummaryGroupDTO

	for rows.Next() {
		var group group.SummaryGroupDTO
		err = rows.Scan(
			&group.GroupID,
			&group.Faculty,
			&group.Program,
			&group.ShortName,
			&group.NumberOfPeople,
			&group.ExistsSchedule)

		if err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}

	return groups, nil
}

func (g *GroupRepository) GetStudentGroupByUserId(ctx context.Context, userID uint64) (group.SummaryGroupDTO, error) {
	sql := `
		SELECT
		    g.group_id,
			f.faculty_name,
			p.program_name,
			g.short_name,
			g.number_of_people,
			g.exists_schedule
		FROM
			public.members as m
		INNER JOIN
			public.groups as g ON g.group_id = m.group_id
		INNER JOIN
			public.faculties as f ON f.faculty_id = g.faculty_id
		INNER JOIN
			public.programs as p ON p.program_id = g.program_id
		WHERE
			user_id = $1
		`

	row := g.pool.QueryRow(ctx, sql, userID)

	var group group.SummaryGroupDTO
	err := row.Scan(
		&group.GroupID,
		&group.Faculty,
		&group.Program,
		&group.ShortName,
		&group.NumberOfPeople,
		&group.ExistsSchedule)

	if err != nil {
		return group, err
	}

	return group, nil
}

func (g *GroupRepository) GetLeaderGroupsByUserId(ctx context.Context, userID uint64) ([]group.DetailsGroupDTO, error) {
	sql := `
		SELECT DISTINCT
			g.group_id,
			g.leader_id,
			f.faculty_name,
			p.program_name,
			g.short_name,
			g.code,
			g.number_of_people,
			g.exists_schedule,
			g.created_at
		FROM
			public.members as m
		INNER JOIN
			public.groups as g ON g.leader_id = m.user_id
		INNER JOIN
			public.faculties as f ON f.faculty_id = g.faculty_id
		INNER JOIN
			public.programs as p ON p.program_id = g.program_id
		WHERE
			user_id = $1
		`

	rows, err := g.pool.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []group.DetailsGroupDTO

	for rows.Next() {
		var group group.DetailsGroupDTO
		err = rows.Scan(
			&group.GroupID,
			&group.LeaderID,
			&group.Faculty,
			&group.Program,
			&group.ShortName,
			&group.Code,
			&group.NumberOfPeople,
			&group.ExistsSchedule,
			&group.CreatedAt)

		if err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}

	return groups, nil
}

func (g *GroupRepository) GetByShortName(ctx context.Context, shortname string) (group.Group, error) {
	sql := `SELECT * FROM public.groups WHERE short_name = $1`

	row := g.pool.QueryRow(ctx, sql, shortname)

	var group group.Group

	err := row.Scan(
		&group.GroupID,
		&group.LeaderID,
		&group.FacultyID,
		&group.ProgramID,
		&group.ShortName,
		&group.Code,
		&group.ExistsSchedule,
		&group.NumberOfPeople,
		&group.CreatedAt)

	if err != nil {
		return group, err
	}

	return group, nil
}

func (g *GroupRepository) GetByCode(ctx context.Context, code string) (group.Group, error) {
	sql := `SELECT * FROM public.groups WHERE code = $1`

	row := g.pool.QueryRow(ctx, sql, code)

	var group group.Group

	err := row.Scan(
		&group.GroupID,
		&group.LeaderID,
		&group.FacultyID,
		&group.ProgramID,
		&group.ShortName,
		&group.Code,
		&group.NumberOfPeople,
		&group.ExistsSchedule,
		&group.CreatedAt)

	if err != nil {
		return group, err
	}

	return group, nil
}

func (g *GroupRepository) GetById(ctx context.Context, groupID uint64) (group.Group, error) {
	sql := `
		SELECT
			*
		FROM
			public.groups
		WHERE group_id = $1
		`

	row := g.pool.QueryRow(ctx, sql, groupID)

	var group group.Group

	err := row.Scan(
		&group.GroupID,
		&group.LeaderID,
		&group.FacultyID,
		&group.ProgramID,
		&group.Code,
		&group.ShortName,
		&group.ExistsSchedule,
		&group.NumberOfPeople,
		&group.CreatedAt)

	if err != nil {
		return group, err
	}

	return group, nil

}
