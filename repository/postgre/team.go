package postgre

import (
	"context"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"gitlab.com/renodesper/spenmo-test/repository"
	"gitlab.com/renodesper/spenmo-test/util/errors"
	"gitlab.com/renodesper/spenmo-test/util/logger"
)

type (
	TeamRepository interface {
		GetAllTeams(ctx context.Context, sortBy string, sort string, skip int, limit int) ([]repository.Team, error)
		GetTeamByID(ctx context.Context, teamID uuid.UUID) (*repository.Team, error)
		GetTeamByName(ctx context.Context, name string) (*repository.Team, error)
		CreateTeam(ctx context.Context, teamPayload *repository.Team) (*repository.Team, error)
		UpdateTeam(ctx context.Context, teamID uuid.UUID, teamPayload map[string]interface{}) (*repository.Team, error)
		DeleteTeam(ctx context.Context, teamID uuid.UUID) (*repository.Team, error)
	}

	TeamRepo struct {
		Log logger.Logger
		Db  *pg.DB
	}
)

var (
	teamTable = "team"
)

// CreateTeamRepository creates team repository
func CreateTeamRepository(log logger.Logger, db *pg.DB) TeamRepository {
	return &TeamRepo{
		Log: log,
		Db:  db,
	}
}

// GetAllTeams ...
func (r *TeamRepo) GetAllTeams(ctx context.Context, sortBy string, sort string, skip int, limit int) ([]repository.Team, error) {
	var teams []repository.Team

	if sortBy == "" {
		sortBy = "created_at"
	}
	if sort == "" {
		sort = "DESC"
	}
	order := fmt.Sprintf("%s %s", sortBy, sort)

	err := r.Db.WithContext(ctx).Model(&teams).Limit(limit).Offset(skip).Order(order).Select()
	if err != nil {
		return nil, errors.FailedTeamsFetch.AppendError(err)
	}

	return teams, nil
}

// GetTeamByID ...
func (r *TeamRepo) GetTeamByID(ctx context.Context, teamID uuid.UUID) (*repository.Team, error) {
	team := repository.Team{}

	sql := r.Db.WithContext(ctx).Model(&team).Where("id = ?", teamID)

	err := sql.Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, errors.FailedNoRows.AppendError(err)
		}

		return nil, errors.FailedTeamFetch.AppendError(err)
	}

	return &team, nil
}

// GetTeamByName ...
func (r *TeamRepo) GetTeamByName(ctx context.Context, name string) (*repository.Team, error) {
	team := repository.Team{}

	sql := r.Db.WithContext(ctx).Model(&team).Where("name = ?", name)

	err := sql.Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, errors.FailedNoRows.AppendError(err)
		}

		return nil, errors.FailedTeamFetch.AppendError(err)
	}

	return &team, nil
}

// CreateTeam ...
func (r *TeamRepo) CreateTeam(ctx context.Context, teamPayload *repository.Team) (*repository.Team, error) {
	var team repository.Team

	_, err := r.Db.WithContext(ctx).Model(teamPayload).Returning("*").Insert(&team)
	if err != nil {
		return nil, errors.FailedTeamCreate.AppendError(err)
	}

	return &team, nil
}

func (r *TeamRepo) UpdateTeam(ctx context.Context, teamID uuid.UUID, teamPayload map[string]interface{}) (*repository.Team, error) {
	teamPayload["updated_at"] = time.Now()

	var team repository.Team
	_, err := r.Db.WithContext(ctx).Model(&teamPayload).Table(teamTable).Where("id = ?", teamID).Returning("*").Update(&team)
	if err != nil {
		return nil, errors.FailedTeamUpdate.AppendError(err)
	}

	return &team, nil
}

func (r *TeamRepo) DeleteTeam(ctx context.Context, teamID uuid.UUID) (*repository.Team, error) {
	teamPayload := map[string]interface{}{
		"is_deleted": true,
		"updated_at": time.Now(),
	}

	var team repository.Team
	_, err := r.Db.WithContext(ctx).Model(&teamPayload).Table(teamTable).Where("id = ?", teamID).Returning("*").Update(&team)
	if err != nil {
		return nil, errors.FailedTeamDelete.AppendError(err)
	}

	return &team, nil
}
