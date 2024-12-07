package group

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/tclutin/classflow-api/internal/domain/edu"
	domainErr "github.com/tclutin/classflow-api/internal/domain/errors"
	"github.com/tclutin/classflow-api/internal/domain/schedule"
	"github.com/tclutin/classflow-api/internal/domain/user"
	"log/slog"
	"time"
)

/*
Нужно решить, что делать с транзакциями и убрать дублирование, думаю в будущем посмотрим, что можно сделать
*/

type UserService interface {
	GetById(ctx context.Context, userID uint64) (user.User, error)
}

type ScheduleService interface {
	GetSchedulesByGroupId(ctx context.Context, filter schedule.FilterDTO, groupID uint64) ([]schedule.DetailsScheduleDTO, error)
}

type EduService interface {
	GetFacultyById(ctx context.Context, facultyID uint64) (edu.Faculty, error)
	GetProgramById(ctx context.Context, programID uint64) (edu.Program, error)
	GetTypeOfSubjectById(ctx context.Context, typeOfSubjectId uint64) (edu.TypeOfSubject, error)
	GetBuildingById(ctx context.Context, buildingID uint64) (edu.Building, error)
}

type UserRepository interface {
	UpdateTx(ctx context.Context, tx pgx.Tx, user user.User) error
}

type ScheduleRepository interface {
	CreateTx(ctx context.Context, tx pgx.Tx, schedule []schedule.Schedule) error
}

type MemberRepository interface {
	DeleteTx(ctx context.Context, tx pgx.Tx, userId uint64) error
	CreateTx(ctx context.Context, tx pgx.Tx, userID uint64, groupId uint64) (uint64, error)
	GetGroupIdByUserId(ctx context.Context, userID uint64) (uint64, error)
}

type Repository interface {
	Create(ctx context.Context, group Group) (uint64, error)
	Update(ctx context.Context, group Group) error
	BeginTx(ctx context.Context) (pgx.Tx, error)
	UpdateTx(ctx context.Context, tx pgx.Tx, group Group) error
	DeleteTx(ctx context.Context, tx pgx.Tx, groupID uint64) error
	GetById(ctx context.Context, groupID uint64) (Group, error)
	GetSummaryGroups(ctx context.Context, filter FilterDTO) ([]SummaryGroupDTO, error)
	GetByShortName(ctx context.Context, shortname string) (Group, error)
	GetDetailsGroupById(ctx context.Context, groupID uint64) (DetailsGroupDTO, error)
}

type Service struct {
	logger          *slog.Logger
	scheduleService ScheduleService
	userService     UserService
	eduService      EduService
	memberRepo      MemberRepository
	scheduleRepo    ScheduleRepository
	userRepo        UserRepository
	repo            Repository
}

func NewService(
	logger *slog.Logger,
	repository Repository,
	memberRepo MemberRepository,
	userRepo UserRepository,
	scheduleService ScheduleService,
	scheduleRepo ScheduleRepository,
	userService UserService,
	eduService EduService,
) *Service {

	return &Service{
		logger:          logger,
		scheduleService: scheduleService,
		scheduleRepo:    scheduleRepo,
		userService:     userService,
		repo:            repository,
		memberRepo:      memberRepo,
		userRepo:        userRepo,
		eduService:      eduService,
	}
}

func (s *Service) Create(ctx context.Context, dto CreateGroupDTO) (uint64, error) {
	_, err := s.GetByShortName(ctx, dto.ShortName)
	if err == nil {
		return 0, domainErr.ErrGroupAlreadyExists
	}

	program, err := s.eduService.GetProgramById(ctx, dto.ProgramID)
	if err != nil {
		return 0, err
	}

	faculty, err := s.eduService.GetFacultyById(ctx, dto.FacultyID)
	if err != nil {
		return 0, err
	}

	if program.FacultyID != faculty.FacultyID {
		return 0, domainErr.ErrFacultyProgramIdMismatch
	}

	entity := Group{
		LeaderID:       nil,
		FacultyID:      dto.FacultyID,
		ProgramID:      dto.ProgramID,
		ShortName:      dto.ShortName,
		NumberOfPeople: 0,
		ExistsSchedule: false,
		CreatedAt:      time.Now(),
	}

	groupID, err := s.repo.Create(ctx, entity)
	if err != nil {
		return 0, fmt.Errorf("error creating group: %w", err)
	}

	return groupID, nil
}

func (s *Service) Delete(ctx context.Context, groupID uint64) error {
	group, err := s.GetById(ctx, groupID)
	if err != nil {
		return err
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			s.logger.Error("Rolling back transaction due to error",
				"error", err,
			)
			tx.Rollback(ctx)
		} else {
			s.logger.Info("Committing transaction")
			tx.Commit(ctx)
		}
	}()

	if group.LeaderID != nil {
		usr, err := s.userService.GetById(ctx, *group.LeaderID)
		if err != nil {
			return err
		}

		usr.Role = user.Student

		if err = s.userRepo.UpdateTx(ctx, tx, usr); err != nil {
			return fmt.Errorf("failed to update user:  %w", err)
		}
	}

	if err = s.repo.DeleteTx(ctx, tx, groupID); err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}

	return nil
}

func (s *Service) Update(ctx context.Context, group Group) error {
	return s.repo.Update(ctx, group)
}

func (s *Service) GetById(ctx context.Context, groupID uint64) (Group, error) {
	group, err := s.repo.GetById(ctx, groupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Group{}, domainErr.ErrGroupNotFound
		}

		return Group{}, fmt.Errorf("failed to get group: %w", err)
	}

	return group, nil
}

func (s *Service) GetByShortName(ctx context.Context, shortname string) (Group, error) {
	group, err := s.repo.GetByShortName(ctx, shortname)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Group{}, domainErr.ErrGroupNotFound
		}

		return Group{}, fmt.Errorf("failed to get group: %w", err)
	}

	return group, nil
}

func (s *Service) GetCurrentGroupByUserID(ctx context.Context, userID uint64) (DetailsGroupDTO, error) {
	groupID, err := s.memberRepo.GetGroupIdByUserId(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DetailsGroupDTO{}, domainErr.ErrMemberNotFound
		}
	}

	currentGroup, err := s.repo.GetDetailsGroupById(ctx, groupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DetailsGroupDTO{}, domainErr.ErrGroupNotFound
		}
	}

	return currentGroup, nil
}

func (s *Service) GetAllGroupsSummary(ctx context.Context, filter FilterDTO) ([]SummaryGroupDTO, error) {
	return s.repo.GetSummaryGroups(ctx, filter)
}

func (s *Service) GetSchedulesByGroupId(ctx context.Context, filter schedule.FilterDTO, groupID uint64) ([]schedule.DetailsScheduleDTO, error) {
	_, err := s.GetById(ctx, groupID)
	if err != nil {
		return nil, err
	}

	schedules, err := s.scheduleService.GetSchedulesByGroupId(ctx, filter, groupID)
	if err != nil {
		return nil, err
	}

	return schedules, nil
}

func (s *Service) UploadSchedule(ctx context.Context, schedule []schedule.Schedule, groupID uint64) error {
	group, err := s.GetById(ctx, groupID)
	if err != nil {
		return err
	}

	if group.ExistsSchedule {
		return domainErr.ErrGroupAlreadyHasSchedule
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			s.logger.Error("Rolling back transaction due to error",
				"error", err,
			)
			tx.Rollback(ctx)
		} else {
			s.logger.Info("Committing transaction")
			tx.Commit(ctx)
		}
	}()

	for _, value := range schedule {
		_, err = s.eduService.GetTypeOfSubjectById(ctx, value.TypeOfSubjectID)
		if err != nil {
			return err
		}

		_, err = s.eduService.GetFacultyById(ctx, value.BuildingsID)
		if err != nil {
			return err
		}

		_, err = s.eduService.GetBuildingById(ctx, value.BuildingsID)
		if err != nil {
			return err
		}
	}

	if err = s.scheduleRepo.CreateTx(ctx, tx, schedule); err != nil {
		return fmt.Errorf("failed to create new schedule: %w", err)
	}

	group.ExistsSchedule = true

	if err = s.repo.UpdateTx(ctx, tx, group); err != nil {
		return fmt.Errorf("failed to update group: %w", err)
	}

	return nil
}

func (s *Service) JoinToGroup(ctx context.Context, userID, groupID uint64) error {
	_, err := s.memberRepo.GetGroupIdByUserId(ctx, userID)
	if err == nil {
		return domainErr.ErrAlreadyInGroup
	}

	group, err := s.GetById(ctx, groupID)
	if err != nil {
		return err
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			s.logger.Error("Rolling back transaction due to error",
				"error", err,
			)
			tx.Rollback(ctx)
		} else {
			s.logger.Info("Committing transaction")
			tx.Commit(ctx)
		}
	}()

	_, err = s.memberRepo.CreateTx(ctx, tx, userID, groupID)
	if err != nil {
		return fmt.Errorf("failed to create member: %w", err)
	}

	group.NumberOfPeople++

	if err = s.repo.UpdateTx(ctx, tx, group); err != nil {
		return fmt.Errorf("failed to update group: %w", err)
	}

	return nil

}

func (s *Service) LeaveFromGroup(ctx context.Context, userID uint64) error {
	groupID, err := s.memberRepo.GetGroupIdByUserId(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainErr.ErrMemberNotFound
		}
	}

	group, err := s.GetById(ctx, groupID)
	if err != nil {
		return err
	}

	usr, err := s.userService.GetById(ctx, userID)
	if err != nil {
		return err
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			s.logger.Error("Rolling back transaction due to error",
				"error", err,
			)
			tx.Rollback(ctx)
		} else {
			s.logger.Info("Committing transaction")
			tx.Commit(ctx)
		}
	}()

	if group.LeaderID != nil && usr.UserID == *group.LeaderID {
		usr.Role = user.Student
		group.LeaderID = nil
		if err = s.userRepo.UpdateTx(ctx, tx, usr); err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	}

	if err = s.memberRepo.DeleteTx(ctx, tx, userID); err != nil {
		return fmt.Errorf("failed to delete member: %w", err)
	}

	group.NumberOfPeople = group.NumberOfPeople - 1

	if err = s.repo.UpdateTx(ctx, tx, group); err != nil {
		return fmt.Errorf("failed to update group: %w", err)
	}

	return nil
}
