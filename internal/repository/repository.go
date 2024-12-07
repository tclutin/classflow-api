package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type Repositories struct {
	User     *UserRepository
	Group    *GroupRepository
	Edu      *EduRepository
	Member   *MemberRepository
	Schedule *ScheduleRepository
}

func NewRepositories(pool *pgxpool.Pool, logger *slog.Logger) *Repositories {
	return &Repositories{
		User:     NewUserRepository(pool, logger),
		Group:    NewGroupRepository(pool),
		Edu:      NewEduRepository(pool, logger),
		Member:   NewMemberRepository(pool),
		Schedule: NewScheduleRepository(pool),
	}
}
