package service

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"gitlab.com/renodesper/spenmo-test/repository"
	"gitlab.com/renodesper/spenmo-test/repository/postgre"
	"gitlab.com/renodesper/spenmo-test/util/errors"
	"gitlab.com/renodesper/spenmo-test/util/logger"
)

type (
	// TeamService ...
	TeamService interface {
		GetAllTeams(ctx context.Context, sortBy string, sort string, skip int, limit int) ([]repository.Team, error)
		GetTeam(ctx context.Context, teamID uuid.UUID) (*repository.Team, error)
		CreateTeam(ctx context.Context, payload *CreateTeamRequest) (*repository.Team, error)
		UpdateTeam(ctx context.Context, teamID uuid.UUID, payload *UpdateTeamRequest) (*repository.Team, error)
		DeleteTeam(ctx context.Context, teamID uuid.UUID) (*repository.Team, error)
	}

	TeamSvc struct {
		Log           logger.Logger
		Team          postgre.TeamRepository
		TeamMemberSvc TeamMemberService
		WalletSvc     WalletService
	}

	CreateTeamRequest struct {
		Name string
	}

	UpdateTeamRequest struct {
		Name string
	}
)

// NewTeamService creates team service
func NewTeamService(log logger.Logger, db *pg.DB) TeamService {
	teamRepo := postgre.CreateTeamRepository(log, db)
	teamMemberSvc := NewTeamMemberService(log, db)
	walletSvc := NewWalletService(log, db)

	return &TeamSvc{
		Log:           log,
		Team:          teamRepo,
		TeamMemberSvc: teamMemberSvc,
		WalletSvc:     walletSvc,
	}
}

func (s *TeamSvc) GetAllTeams(ctx context.Context, sortBy string, sort string, skip int, limit int) ([]repository.Team, error) {
	var teams []repository.Team

	teams, err := s.Team.GetAllTeams(ctx, sortBy, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	if teams == nil {
		teams = []repository.Team{}
	}

	return teams, nil
}

func (s *TeamSvc) GetTeam(ctx context.Context, teamID uuid.UUID) (*repository.Team, error) {
	team, err := s.Team.GetTeamByID(ctx, teamID)
	return team, err
}

func (s *TeamSvc) CreateTeam(ctx context.Context, payload *CreateTeamRequest) (*repository.Team, error) {
	if payload.Name != "" {
		team, _ := s.Team.GetTeamByName(ctx, payload.Name)

		if team != nil {
			return nil, errors.FailedTeamExist
		}
	}

	ID := uuid.New()

	teamPayload := repository.Team{
		ID:   ID,
		Name: payload.Name,
	}
	team, err := s.Team.CreateTeam(ctx, &teamPayload)
	return team, err
}

func (s *TeamSvc) UpdateTeam(ctx context.Context, teamID uuid.UUID, payload *UpdateTeamRequest) (*repository.Team, error) {
	team, _ := s.Team.GetTeamByID(ctx, teamID)

	if team != nil && teamID != team.ID {
		return nil, errors.FailedTeamExist
	}

	var teamPayload = make(map[string]interface{})

	if payload.Name != "" {
		teamPayload["name"] = payload.Name
	}

	team, err := s.Team.UpdateTeam(ctx, teamID, teamPayload)
	return team, err
}

func (s *TeamSvc) DeleteTeam(ctx context.Context, teamID uuid.UUID) (*repository.Team, error) {
	team, err := s.Team.DeleteTeam(ctx, teamID)
	if err != nil {
		return nil, err
	}

	_, err = s.TeamMemberSvc.DeleteTeamMembersByTeamID(ctx, teamID)
	if err != nil {
		s.Log.Warn(err)
	}

	_, err = s.WalletSvc.DeleteWalletsByTeamID(ctx, teamID)
	if err != nil {
		s.Log.Warn(err)
	}

	return team, nil
}
