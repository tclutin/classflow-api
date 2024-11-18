package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repositories struct {
	User   *UserRepository
	Group  *GroupRepository
	Edu    *EduRepository
	Member *MemberRepository
}

func NewRepositories(pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		User:   NewUserRepository(pool),
		Group:  NewGroupRepository(pool),
		Edu:    NewEduRepository(pool),
		Member: NewMemberRepository(pool),
	}
}
