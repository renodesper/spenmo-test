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
	// UserService ...
	UserService interface {
		GetAllUsers(ctx context.Context, sortBy string, sort string, skip int, limit int) ([]repository.User, error)
		GetUser(ctx context.Context, userID uuid.UUID) (*repository.User, error)
		CreateUser(ctx context.Context, payload *CreateUserRequest) (*repository.User, error)
		UpdateUser(ctx context.Context, userID uuid.UUID, payload *UpdateUserRequest) (*repository.User, error)
		DeleteUser(ctx context.Context, userID uuid.UUID) (*repository.User, error)
	}

	UserSvc struct {
		Log           logger.Logger
		User          postgre.UserRepository
		TeamMemberSvc TeamMemberService
		WalletSvc     WalletService
	}

	CreateUserRequest struct {
		Email string
		Name  string
	}

	UpdateUserRequest struct {
		Email string
		Name  string
	}
)

// NewUserService creates user service
func NewUserService(log logger.Logger, db *pg.DB) UserService {
	userRepo := postgre.CreateUserRepository(log, db)
	teamMemberSvc := NewTeamMemberService(log, db)
	walletSvc := NewWalletService(log, db)

	return &UserSvc{
		Log:           log,
		User:          userRepo,
		TeamMemberSvc: teamMemberSvc,
		WalletSvc:     walletSvc,
	}
}

func (s *UserSvc) GetAllUsers(ctx context.Context, sortBy string, sort string, skip int, limit int) ([]repository.User, error) {
	var users []repository.User

	users, err := s.User.GetAllUsers(ctx, sortBy, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	if users == nil {
		users = []repository.User{}
	}

	return users, nil
}

func (s *UserSvc) GetUser(ctx context.Context, userID uuid.UUID) (*repository.User, error) {
	user, err := s.User.GetUserByID(ctx, userID)
	return user, err
}

func (s *UserSvc) CreateUser(ctx context.Context, payload *CreateUserRequest) (*repository.User, error) {
	if payload.Email != "" {
		user, _ := s.User.GetUserByEmail(ctx, payload.Email)

		if user != nil {
			return nil, errors.FailedEmailExist
		}
	}

	ID := uuid.New()

	userPayload := repository.User{
		ID:    ID,
		Email: payload.Email,
		Name:  payload.Name,
	}
	user, err := s.User.CreateUser(ctx, &userPayload)
	return user, err
}

func (s *UserSvc) UpdateUser(ctx context.Context, userID uuid.UUID, payload *UpdateUserRequest) (*repository.User, error) {
	if payload.Email != "" {
		user, _ := s.User.GetUserByEmail(ctx, payload.Email)

		if user != nil && userID != user.ID {
			return nil, errors.FailedEmailExist
		}
	}

	var userPayload = make(map[string]interface{})

	if payload.Email != "" {
		userPayload["email"] = payload.Email
	}

	if payload.Name != "" {
		userPayload["name"] = payload.Name
	}

	user, err := s.User.UpdateUser(ctx, userID, userPayload)
	return user, err
}

func (s *UserSvc) DeleteUser(ctx context.Context, userID uuid.UUID) (*repository.User, error) {
	user, err := s.User.DeleteUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	_, err = s.TeamMemberSvc.DeleteTeamMembersByUserID(ctx, userID)
	if err != nil {
		s.Log.Warn(err)
	}

	_, err = s.WalletSvc.DeleteWalletsByUserID(ctx, userID)
	if err != nil {
		s.Log.Warn(err)
	}

	return user, nil
}
