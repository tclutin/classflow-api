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
	"github.com/tclutin/classflow-api/pkg/hash"
	"time"
)

type UserService interface {
	GetById(ctx context.Context, userID uint64) (user.User, error)
}

type ScheduleService interface {
	Create(ctx context.Context, schedule []schedule.Schedule) error
	GetAllSchedulesByGroupId(ctx context.Context, groupID uint64) ([]schedule.DetailsScheduleDTO, error)
}

type EduService interface {
	GetFacultyById(ctx context.Context, facultyID uint64) (edu.Faculty, error)
	GetProgramById(ctx context.Context, programID uint64) (edu.Program, error)
	GetTypeOfSubjectById(ctx context.Context, typeOfSubjectId uint64) (edu.TypeOfSubject, error)
	GetBuildingById(ctx context.Context, buildingID uint64) (edu.Building, error)
}

type MemberRepository interface {
	Delete(ctx context.Context, userId uint64) error
	Create(ctx context.Context, userID uint64, groupId uint64) (uint64, error)
	GetGroupIdByUserId(ctx context.Context, userID uint64) (uint64, error)
}

type Repository interface {
	Create(ctx context.Context, group Group) (uint64, error)
	Update(ctx context.Context, group Group) error
	GetById(ctx context.Context, groupID uint64) (Group, error)
	GetSummaryGroups(ctx context.Context, filter FilterDTO) ([]SummaryGroupDTO, error)
	GetByShortName(ctx context.Context, shortname string) (Group, error)
	GetStudentGroupByUserId(ctx context.Context, userID uint64) (SummaryGroupDTO, error)
	GetLeaderGroupsByUserId(ctx context.Context, userID uint64) ([]DetailsGroupDTO, error)
}

type Service struct {
	scheduleService ScheduleService
	userService     UserService
	eduService      EduService
	memberRepo      MemberRepository
	repo            Repository
}

func NewService(
	repository Repository,
	memberRepo MemberRepository,
	scheduleService ScheduleService,
	userService UserService,
	eduService EduService,
) *Service {

	return &Service{
		scheduleService: scheduleService,
		userService:     userService,
		repo:            repository,
		memberRepo:      memberRepo,
		eduService:      eduService,
	}
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

	//TODO TX manger needs
	if err = s.memberRepo.Delete(ctx, userID); err != nil {
		return err
	}

	group.NumberOfPeople = group.NumberOfPeople - 1

	return s.repo.Update(ctx, group)
}

func (s *Service) GetAllSchedulesByGroupIdAndUserId(ctx context.Context, groupID, userID uint64) ([]schedule.DetailsScheduleDTO, error) {
	group, err := s.GetById(ctx, groupID)
	if err != nil {
		return nil, err
	}

	groupID, err = s.memberRepo.GetGroupIdByUserId(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.ErrYouArentMember
		}
	}

	if group.GroupID != groupID {
		return nil, domainErr.ErrYouArentMember
	}

	schedules, err := s.scheduleService.GetAllSchedulesByGroupId(ctx, groupID)
	if err != nil {
		return nil, err
	}

	return schedules, nil
}

// TODO: need tx manager on service layer
func (s *Service) UploadSchedule(ctx context.Context, schedule []schedule.Schedule, groupID, userID uint64) error {
	group, err := s.GetById(ctx, groupID)
	if err != nil {
		return err
	}

	if group.LeaderID != userID {
		return domainErr.ErrThisGroupDoesNotBelongToYou
	}

	if group.ExistsSchedule {
		return domainErr.ErrGroupAlreadyHasSchedule
	}

	for _, value := range schedule {
		_, err = s.eduService.GetTypeOfSubjectById(ctx, value.TypeOfSubjectID)
		if err != nil {
			return err
		}

		_, err = s.eduService.GetFacultyById(ctx, value.BuildingsID)
		if err != nil {
			return err
		}
	}

	if err = s.scheduleService.Create(ctx, schedule); err != nil {
		return fmt.Errorf("cannot create new schedule: %w", err)
	}

	group.ExistsSchedule = true

	if err = s.Update(ctx, group); err != nil {
		return fmt.Errorf("cannot update group: %w", err)
	}

	return nil
}

// TODO: tx and normiы
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

	code, err := s.GenCode(4)
	if err != nil {
		return 0, err
	}

	entity := Group{
		LeaderID:       dto.LeaderID,
		FacultyID:      dto.FacultyID,
		ProgramID:      dto.ProgramID,
		ShortName:      dto.ShortName,
		Code:           code,
		NumberOfPeople: 1,
		ExistsSchedule: false,
		CreatedAt:      time.Now(),
	}

	groupID, err := s.repo.Create(ctx, entity)
	if err != nil {
		return 0, fmt.Errorf("error creating group: %w", err)
	}

	_, err = s.memberRepo.Create(ctx, entity.LeaderID, groupID)
	if err != nil {
		return 0, err
	}

	return groupID, nil
}

func (s *Service) Update(ctx context.Context, group Group) error {
	return s.repo.Update(ctx, group)
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

func (s *Service) GetStudentGroupByUserId(ctx context.Context, userID uint64) (SummaryGroupDTO, error) {
	group, err := s.repo.GetStudentGroupByUserId(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return SummaryGroupDTO{}, domainErr.ErrGroupNotFound
		}

		return SummaryGroupDTO{}, fmt.Errorf("failed to get student group: %w", err)
	}

	return group, nil
}

func (s *Service) GetLeaderGroupsByUserId(ctx context.Context, userID uint64) ([]DetailsGroupDTO, error) {
	return s.repo.GetLeaderGroupsByUserId(ctx, userID)
}

func (s *Service) GetAllGroupsSummary(ctx context.Context, filter FilterDTO) ([]SummaryGroupDTO, error) {
	return s.repo.GetSummaryGroups(ctx, filter)
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

func (s *Service) JoinToGroup(ctx context.Context, code string, userID, groupID uint64) error {
	_, err := s.memberRepo.GetGroupIdByUserId(ctx, userID)
	if err == nil {
		return domainErr.ErrAlreadyInGroup
	}

	group, err := s.repo.GetById(ctx, groupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainErr.ErrGroupNotFound
		}

		return fmt.Errorf("failed to get group: %w", err)
	}

	if group.Code != code {
		return domainErr.ErrWrongGroupCode
	}

	_, err = s.memberRepo.Create(ctx, userID, groupID)
	if err != nil {
		return fmt.Errorf("failed to create memberp: %w", err)
	}

	group.NumberOfPeople++

	return s.repo.Update(ctx, group)

}

func (s *Service) GenCode(size int64) (string, error) {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	alias := make([]rune, size)

	for i := range alias {
		rnd, err := hash.NewCryptoRand(int64(len(chars)))
		if err != nil {
			return "", err
		}
		alias[i] = chars[rnd]
	}

	return string(alias), nil
}
