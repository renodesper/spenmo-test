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
	TeamMemberRepository interface {
		GetTeamMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) (*repository.TeamMember, error)
		GetTeamMembers(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, sortBy string, sort string, skip int, limit int) ([]repository.TeamMember, error)
		CreateTeamMember(ctx context.Context, teamMemberPayload *repository.TeamMember) (*repository.TeamMember, error)
		DeleteTeamMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) (*repository.TeamMember, error)
		DeleteTeamMembersByTeamID(ctx context.Context, teamID uuid.UUID) ([]repository.TeamMember, error)
		DeleteTeamMembersByUserID(ctx context.Context, userID uuid.UUID) ([]repository.TeamMember, error)
	}

	TeamMemberRepo struct {
		Log logger.Logger
		Db  *pg.DB
	}
)

var (
	teamMemberTable = "team_member"
)

// CreateTeamMemberRepository creates teamMember repository
func CreateTeamMemberRepository(log logger.Logger, db *pg.DB) TeamMemberRepository {
	return &TeamMemberRepo{
		Log: log,
		Db:  db,
	}
}

func (r *TeamMemberRepo) GetTeamMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) (*repository.TeamMember, error) {
	teamMember := repository.TeamMember{}

	sql := r.Db.WithContext(ctx).Model(&teamMember).Where("team_id = ?", teamID).Where("user_id = ?", userID)

	err := sql.Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, errors.FailedNoRows.AppendError(err)
		}

		return nil, errors.FailedTeamMemberFetch.AppendError(err)
	}

	return &teamMember, nil
}

func (r *TeamMemberRepo) GetTeamMembers(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, sortBy string, sort string, skip int, limit int) ([]repository.TeamMember, error) {
	teamMembers := []repository.TeamMember{}

	if sortBy == "" {
		sortBy = "created_at"
	}
	if sort == "" {
		sort = "DESC"
	}
	order := fmt.Sprintf("%s %s", sortBy, sort)

	sql := r.Db.WithContext(ctx).Model(&teamMembers)

	if teamID != uuid.Nil {
		sql = sql.Where("team_id = ?", teamID)
	}

	if userID != uuid.Nil {
		sql = sql.Where("user_id = ?", userID)
	}

	err := sql.Limit(limit).Offset(skip).Order(order).Select()
	if err != nil {
		return nil, errors.FailedTeamMembersFetch.AppendError(err)
	}

	return teamMembers, nil
}

// CreateTeamMember ...
func (r *TeamMemberRepo) CreateTeamMember(ctx context.Context, teamMemberPayload *repository.TeamMember) (*repository.TeamMember, error) {
	var teamMember repository.TeamMember

	_, err := r.Db.WithContext(ctx).Model(teamMemberPayload).Returning("*").Insert(&teamMember)
	if err != nil {
		return nil, errors.FailedTeamMemberCreate.AppendError(err)
	}

	return &teamMember, nil
}

func (r *TeamMemberRepo) DeleteTeamMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) (*repository.TeamMember, error) {
	teamMemberPayload := map[string]interface{}{
		"is_deleted": true,
		"updated_at": time.Now(),
	}

	var teamMember repository.TeamMember
	_, err := r.Db.WithContext(ctx).Model(&teamMemberPayload).Table(teamMemberTable).Where("team_id = ?", teamID).Where("user_id = ?", userID).Returning("*").Update(&teamMember)
	if err != nil {
		return nil, errors.FailedTeamMemberDelete.AppendError(err)
	}

	return &teamMember, nil
}

func (r *TeamMemberRepo) DeleteTeamMembersByTeamID(ctx context.Context, teamID uuid.UUID) ([]repository.TeamMember, error) {
	teamMemberPayload := map[string]interface{}{
		"is_deleted": true,
		"updated_at": time.Now(),
	}

	var teamMembers []repository.TeamMember
	_, err := r.Db.WithContext(ctx).Model(&teamMemberPayload).Table(teamMemberTable).Where("team_id = ?", teamID).Returning("*").Update(&teamMembers)
	if err != nil {
		return nil, errors.FailedTeamMemberDelete.AppendError(err)
	}

	return teamMembers, nil
}

func (r *TeamMemberRepo) DeleteTeamMembersByUserID(ctx context.Context, userID uuid.UUID) ([]repository.TeamMember, error) {
	teamMemberPayload := map[string]interface{}{
		"is_deleted": true,
		"updated_at": time.Now(),
	}

	var teamMembers []repository.TeamMember
	_, err := r.Db.WithContext(ctx).Model(&teamMemberPayload).Table(teamMemberTable).Where("user_id = ?", userID).Returning("*").Update(&teamMembers)
	if err != nil {
		return nil, errors.FailedTeamMemberDelete.AppendError(err)
	}

	return teamMembers, nil
}
