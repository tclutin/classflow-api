package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tclutin/classflow-api/internal/domain/group"
	"log/slog"
	"strings"
)

type GroupRepository struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewGroupRepository(pool *pgxpool.Pool, logger *slog.Logger) *GroupRepository {
	return &GroupRepository{
		pool:   pool,
		logger: logger,
	}
}

func (g *GroupRepository) Create(ctx context.Context, group group.Group) (uint64, error) {
	sql := `
	INSERT INTO public.groups
    (leader_id, faculty_id, program_id, short_name, exists_schedule, number_of_people, created_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING group_id;`

	row := g.pool.QueryRow(
		ctx,
		sql,
		group.LeaderID,
		group.FacultyID,
		group.ProgramID,
		group.ShortName,
		group.ExistsSchedule,
		group.NumberOfPeople,
		group.CreatedAt)

	var groupId uint64

	if err := row.Scan(&groupId); err != nil {
		g.logger.Error("Failed to create group",
			"error", err,
			"faculty_id", group.FacultyID,
			"program_id", group.ProgramID,
			"short_name", group.ShortName,
		)
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
			number_of_people = $5,
			exists_schedule = $6,
			created_at = $7
		WHERE
		    group_id = $8
		`

	_, err := g.pool.Exec(
		ctx,
		sql,
		group.LeaderID,
		group.FacultyID,
		group.ProgramID,
		group.ShortName,
		group.NumberOfPeople,
		group.ExistsSchedule,
		group.CreatedAt,
		group.GroupID)

	if err != nil {
		g.logger.Error("Failed to update group",
			"error", err,
			"group_id", group.GroupID,
		)
		return err
	}

	return nil
}

func (g *GroupRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return g.pool.Begin(ctx)
}

func (g *GroupRepository) UpdateTx(ctx context.Context, tx pgx.Tx, group group.Group) error {
	sql := `
		UPDATE
			public.groups
		SET
			leader_id = $1,
			faculty_id = $2,
			program_id = $3,
			short_name = $4,
			number_of_people = $5,
			exists_schedule = $6,
			created_at = $7
		WHERE
		    group_id = $8
		`

	_, err := tx.Exec(
		ctx,
		sql,
		group.LeaderID,
		group.FacultyID,
		group.ProgramID,
		group.ShortName,
		group.NumberOfPeople,
		group.ExistsSchedule,
		group.CreatedAt,
		group.GroupID)

	if err != nil {
		g.logger.Error("Failed to update group",
			"error", err,
			"group_id", group.GroupID,
		)
		return err
	}

	return nil
}

func (g *GroupRepository) DeleteTx(ctx context.Context, tx pgx.Tx, groupID uint64) error {
	sql := `DELETE FROM public.groups WHERE group_id = $1`

	_, err := tx.Exec(ctx, sql, groupID)

	if err != nil {
		g.logger.Error("Failed to delete group",
			"error", err,
			"group_id", groupID,
		)
		return err
	}

	return nil
}

// GetSummaryGroups TODO: needs sql builder/string builder
func (g *GroupRepository) GetSummaryGroups(ctx context.Context, filter group.FilterDTO) ([]group.SummaryGroupDTO, error) {
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

	var conditions []string
	var args []interface{}
	var argCount int

	if filter.Faculty != "" {
		conditions = append(conditions, fmt.Sprintf("f.faculty_name = $%d", argCount+1))
		args = append(args, filter.Faculty)
		argCount++
	}

	if filter.Program != "" {
		conditions = append(conditions, fmt.Sprintf("p.program_name = $%d", argCount+1))
		args = append(args, filter.Program)
		argCount++
	}

	if len(conditions) > 0 {
		sql += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := g.pool.Query(ctx, sql, args...)
	if err != nil {
		g.logger.Error("Failed to get summary groups",
			"error", err,
			"args", args,
		)
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
			g.logger.Error("Failed to scan row in GetSummaryGroups",
				"error", err,
			)
			return nil, err
		}

		groups = append(groups, group)
	}

	return groups, nil
}

func (g *GroupRepository) GetDetailsGroupById(ctx context.Context, groupID uint64) (group.DetailsGroupDTO, error) {
	sql := `
		SELECT 
			g.group_id,
			g.leader_id,
			f.faculty_name,
			p.program_name,
			g.short_name,
			g.number_of_people,
			g.exists_schedule,
			g.created_at
		FROM
			public.groups AS g
		INNER JOIN
			public.programs AS p ON g.program_id = p.program_id
		INNER JOIN
			public.faculties AS f ON g.faculty_id = f.faculty_id
		WHERE
			group_id = $1
		`

	row := g.pool.QueryRow(ctx, sql, groupID)

	var group group.DetailsGroupDTO
	err := row.Scan(
		&group.GroupID,
		&group.LeaderID,
		&group.Faculty,
		&group.Program,
		&group.ShortName,
		&group.NumberOfPeople,
		&group.ExistsSchedule,
		&group.CreatedAt)

	// TODO: если нет такой записи, то сделать варн
	if err != nil {
		g.logger.Error("Failed to get detail group with group id",
			"error", err,
			"group_id", groupID,
		)
		return group, err
	}

	return group, nil
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
		&group.ExistsSchedule,
		&group.NumberOfPeople,
		&group.CreatedAt)

	if err != nil {
		g.logger.Error("Failed to get group by shortname",
			"error", err,
			"short_name", shortname,
		)
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
		&group.ShortName,
		&group.ExistsSchedule,
		&group.NumberOfPeople,
		&group.CreatedAt)

	if err != nil {
		g.logger.Error("Failed to get group by id",
			"error", err,
			"group_id", groupID,
		)
		return group, err
	}

	return group, nil

}
