package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tclutin/classflow-api/internal/domain/schedule"
)

type ScheduleRepository struct {
	pool *pgxpool.Pool
}

func NewScheduleRepository(pool *pgxpool.Pool) *ScheduleRepository {
	return &ScheduleRepository{pool}
}

func (s *ScheduleRepository) Create(ctx context.Context, schedule []schedule.Schedule) error {
	sql := `
		INSERT INTO public.schedule
		(group_id, buildings_id, type_of_subject_id, subject_name, teacher, room, is_even, day_of_week, start_time, end_time, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, value := range schedule {
		_, err = tx.Exec(
			ctx,
			sql,
			value.GroupID,
			value.BuildingsID,
			value.TypeOfSubjectID,
			value.SubjectName,
			value.Teacher,
			value.Room,
			value.IsEven,
			value.DayOfWeek,
			value.StartTime,
			value.EndTime,
			value.CreatedAt)

		if err != nil {
			return err
		}

	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (s *ScheduleRepository) GetAllSchedulesByGroupId(ctx context.Context, filter schedule.FilterDTO, groupID uint64) ([]schedule.DetailsScheduleDTO, error) {
	sql := `
		SELECT
			t.name,
			s.subject_name,
			s.teacher,
			s.room,
			s.is_even,
			s.day_of_week,
			s.start_time,
			s.end_time,
			b.buildings_id,
			b.name,
			b.latitude,
			b.longitude,
			b.address
		FROM
			public.schedule as s
		INNER JOIN
			public.type_of_subject as t ON s.type_of_subject_id = t.type_of_subject_id
		INNER JOIN
			public.buildings as b ON s.buildings_id = b.buildings_id
		WHERE
			group_id = $1
		`

	if filter.IsEven == "true" {
		sql += " AND s.is_even = true"
	}

	if filter.IsEven == "false" {
		sql += " AND s.is_even = false"
	}

	rows, err := s.pool.Query(ctx, sql, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []schedule.DetailsScheduleDTO

	for rows.Next() {
		var schedule schedule.DetailsScheduleDTO
		err = rows.Scan(
			&schedule.Type,
			&schedule.SubjectName,
			&schedule.Teacher,
			&schedule.Room,
			&schedule.IsEven,
			&schedule.DayOfWeek,
			&schedule.StartTime,
			&schedule.EndTime,
			&schedule.Building.BuildingID,
			&schedule.Building.Name,
			&schedule.Building.Latitude,
			&schedule.Building.Longitude,
			&schedule.Building.Address)

		if err != nil {
			return nil, err
		}

		schedules = append(schedules, schedule)
	}

	return schedules, nil
}
