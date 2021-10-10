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
	UserRepository interface {
		GetAllUsers(ctx context.Context, sortBy string, sort string, skip int, limit int) ([]repository.User, error)
		GetUserByID(ctx context.Context, userID uuid.UUID) (*repository.User, error)
		GetUserByEmail(ctx context.Context, email string) (*repository.User, error)
		CreateUser(ctx context.Context, userPayload *repository.User) (*repository.User, error)
		UpdateUser(ctx context.Context, userID uuid.UUID, userPayload map[string]interface{}) (*repository.User, error)
		DeleteUser(ctx context.Context, userID uuid.UUID) (*repository.User, error)
	}

	UserRepo struct {
		Log logger.Logger
		Db  *pg.DB
	}
)

var (
	userTable = "user"
)

// CreateUserRepository creates user repository
func CreateUserRepository(log logger.Logger, db *pg.DB) UserRepository {
	return &UserRepo{
		Log: log,
		Db:  db,
	}
}

// GetAllUsers ...
func (r *UserRepo) GetAllUsers(ctx context.Context, sortBy string, sort string, skip int, limit int) ([]repository.User, error) {
	var users []repository.User

	if sortBy == "" {
		sortBy = "created_at"
	}
	if sort == "" {
		sort = "DESC"
	}
	order := fmt.Sprintf("%s %s", sortBy, sort)

	err := r.Db.WithContext(ctx).Model(&users).Limit(limit).Offset(skip).Order(order).Select()
	if err != nil {
		return nil, errors.FailedUsersFetch.AppendError(err)
	}

	return users, nil
}

// GetUserByID ...
func (r *UserRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (*repository.User, error) {
	user := repository.User{}

	sql := r.Db.WithContext(ctx).Model(&user).Where("id = ?", userID)

	err := sql.Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, errors.FailedNoRows.AppendError(err)
		}

		return nil, errors.FailedUserFetch.AppendError(err)
	}

	return &user, nil
}

// GetUserByEmail ...
func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*repository.User, error) {
	user := repository.User{}

	sql := r.Db.WithContext(ctx).Model(&user).Where("email = ?", email)

	err := sql.Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, errors.FailedNoRows.AppendError(err)
		}

		return nil, errors.FailedUserFetch.AppendError(err)
	}

	return &user, nil
}

// CreateUser ...
func (r *UserRepo) CreateUser(ctx context.Context, userPayload *repository.User) (*repository.User, error) {
	var user repository.User

	_, err := r.Db.WithContext(ctx).Model(userPayload).Returning("*").Insert(&user)
	if err != nil {
		return nil, errors.FailedUserCreate.AppendError(err)
	}

	return &user, nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, userID uuid.UUID, userPayload map[string]interface{}) (*repository.User, error) {
	userPayload["updated_at"] = time.Now()

	var user repository.User
	_, err := r.Db.WithContext(ctx).Model(&userPayload).Table(userTable).Where("id = ?", userID).Returning("*").Update(&user)
	if err != nil {
		return nil, errors.FailedUserUpdate.AppendError(err)
	}

	return &user, nil
}

func (r *UserRepo) DeleteUser(ctx context.Context, userID uuid.UUID) (*repository.User, error) {
	userPayload := map[string]interface{}{
		"is_deleted": true,
		"updated_at": time.Now(),
	}

	var user repository.User
	_, err := r.Db.WithContext(ctx).Model(&userPayload).Table(userTable).Where("id = ?", userID).Returning("*").Update(&user)
	if err != nil {
		return nil, errors.FailedUserDelete.AppendError(err)
	}

	return &user, nil
}
