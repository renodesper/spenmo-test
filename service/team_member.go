package service

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"gitlab.com/renodesper/spenmo-test/repository"
	"gitlab.com/renodesper/spenmo-test/repository/postgre"
	e "gitlab.com/renodesper/spenmo-test/util/error"
	"gitlab.com/renodesper/spenmo-test/util/errors"
	"gitlab.com/renodesper/spenmo-test/util/logger"
)

type (
	// TeamMemberService ...
	TeamMemberService interface {
		GetTeamMembers(ctx context.Context, payload *GetTeamMembersRequest) ([]repository.TeamMember, error)
		CreateTeamMember(ctx context.Context, payload *CreateTeamMemberRequest) (*repository.TeamMember, error)
		DeleteTeamMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) (*repository.TeamMember, error)
		DeleteTeamMembersByTeamID(ctx context.Context, teamID uuid.UUID) ([]repository.TeamMember, error)
		DeleteTeamMembersByUserID(ctx context.Context, userID uuid.UUID) ([]repository.TeamMember, error)
	}

	TeamMemberSvc struct {
		Log        logger.Logger
		TeamMember postgre.TeamMemberRepository
	}

	CreateTeamMemberRequest struct {
		TeamID uuid.UUID
		UserID uuid.UUID
	}

	GetTeamMembersRequest struct {
		TeamID uuid.UUID
		UserID uuid.UUID
		SortBy string
		Sort   string
		Skip   int
		Limit  int
	}
)

// NewTeamMemberService creates team service
func NewTeamMemberService(log logger.Logger, db *pg.DB) TeamMemberService {
	teamRepo := postgre.CreateTeamMemberRepository(log, db)

	return &TeamMemberSvc{
		Log:        log,
		TeamMember: teamRepo,
	}
}

func (s *TeamMemberSvc) GetTeamMembers(ctx context.Context, payload *GetTeamMembersRequest) ([]repository.TeamMember, error) {
	teamMembers, err := s.TeamMember.GetTeamMembers(ctx, payload.TeamID, payload.UserID, payload.SortBy, payload.Sort, payload.Skip, payload.Limit)
	return teamMembers, err
}

func (s *TeamMemberSvc) CreateTeamMember(ctx context.Context, teamMemberPayload *CreateTeamMemberRequest) (*repository.TeamMember, error) {
	if teamMemberPayload.TeamID == uuid.Nil {
		return nil, errors.MissingTeamID
	}

	if teamMemberPayload.UserID == uuid.Nil {
		return nil, errors.MissingUserID
	}

	existingTeamMember, err := s.TeamMember.GetTeamMember(ctx, teamMemberPayload.TeamID, teamMemberPayload.UserID)
	if err != nil {
		if er := err.(e.Error); er.Code != errors.FailedNoRows.Code {
			return nil, err
		}
	}
	if existingTeamMember != nil {
		return nil, errors.FailedTeamMemberExist
	}

	ID := uuid.New()

	teamPayload := repository.TeamMember{
		ID:     ID,
		TeamID: teamMemberPayload.TeamID,
		UserID: teamMemberPayload.UserID,
	}
	teamMember, err := s.TeamMember.CreateTeamMember(ctx, &teamPayload)
	return teamMember, err
}

func (s *TeamMemberSvc) DeleteTeamMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) (*repository.TeamMember, error) {
	teamMember, err := s.TeamMember.DeleteTeamMember(ctx, teamID, userID)
	return teamMember, err
}

func (s *TeamMemberSvc) DeleteTeamMembersByTeamID(ctx context.Context, teamID uuid.UUID) ([]repository.TeamMember, error) {
	teamMembers, err := s.TeamMember.DeleteTeamMembersByTeamID(ctx, teamID)
	return teamMembers, err
}

func (s *TeamMemberSvc) DeleteTeamMembersByUserID(ctx context.Context, userID uuid.UUID) ([]repository.TeamMember, error) {
	teamMembers, err := s.TeamMember.DeleteTeamMembersByUserID(ctx, userID)
	return teamMembers, err
}
