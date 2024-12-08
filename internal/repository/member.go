package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type MemberRepository struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewMemberRepository(pool *pgxpool.Pool, logger *slog.Logger) *MemberRepository {
	return &MemberRepository{
		pool:   pool,
		logger: logger,
	}
}

func (m *MemberRepository) CreateTx(ctx context.Context, tx pgx.Tx, userID uint64, groupId uint64) (uint64, error) {
	sql := `INSERT INTO public.members (user_id, group_id) VALUES ($1, $2) RETURNING member_id`

	row := tx.QueryRow(ctx, sql, userID, groupId)

	var memberID uint64

	if err := row.Scan(&memberID); err != nil {
		m.logger.Error("Failed to create member",
			"error", err,
			"userID", userID,
			"groupId", groupId,
		)
		return 0, err
	}

	return memberID, nil
}

func (m *MemberRepository) DeleteTx(ctx context.Context, tx pgx.Tx, userId uint64) error {
	sql := `DELETE FROM public.members WHERE user_id = $1`

	_, err := tx.Exec(ctx, sql, userId)

	if err != nil {
		m.logger.Error("Failed to delete member",
			"error", err,
			"user_id", userId,
		)
		return err
	}

	return nil
}

func (m *MemberRepository) GetGroupIdByUserId(ctx context.Context, userID uint64) (uint64, error) {
	sql := `SELECT group_id FROM public.members WHERE user_id = $1`

	row := m.pool.QueryRow(ctx, sql, userID)

	var memberID uint64
	if err := row.Scan(&memberID); err != nil {
		m.logger.Error("Failed to get group ID for user",
			"error", err,
			"userID", userID,
		)
		return 0, err
	}

	return memberID, nil
}
