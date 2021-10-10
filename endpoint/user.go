package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"gitlab.com/renodesper/spenmo-test/service"
)

type (
	CreateUserRequest struct {
		Email string `json:"email" validate:"required,email"`
		Name  string `json:"name" validate:"required"`
	}

	GetAllUsersRequest struct {
		SortBy string `json:"sortBy"`
		Sort   string `json:"sort"`
		Skip   int    `json:"skip"`
		Limit  int    `json:"limit"`
	}

	GetUserRequest struct {
		ID uuid.UUID `json:"id" validate:"required"`
	}

	UpdateUserRequest struct {
		ID    uuid.UUID `json:"id" validate:"required"`
		Email string
		Name  string
	}

	DeleteUserRequest struct {
		ID uuid.UUID `json:"id" validate:"required"`
	}
)

func MakeCreateUserEndpoint(userSvc service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateUserRequest)

		userReq := service.CreateUserRequest{
			Email: req.Email,
			Name:  req.Name,
		}

		user, err := userSvc.CreateUser(ctx, &userReq)
		return user, err
	}
}

func MakeGetAllUsersEndpoint(userSvc service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetAllUsersRequest)

		users, err := userSvc.GetAllUsers(ctx, req.SortBy, req.Sort, req.Skip, req.Limit)
		if err != nil {
			return nil, err
		}

		response = map[string]interface{}{
			"users": users,
			"pagination": map[string]interface{}{
				"sortBy": req.SortBy,
				"sort":   req.Sort,
				"skip":   req.Skip,
				"limit":  req.Limit,
			},
		}

		return response, nil
	}
}

func MakeGetUserByIDEndpoint(userSvc service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetUserRequest)

		user, err := userSvc.GetUser(ctx, req.ID)
		return user, err
	}
}

func MakeUpdateUserEndpoint(userSvc service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateUserRequest)

		userPayload := service.UpdateUserRequest{
			Email: req.Email,
			Name:  req.Name,
		}

		user, err := userSvc.UpdateUser(ctx, req.ID, &userPayload)
		return user, err
	}
}

func MakeDeleteUserByIDEndpoint(userSvc service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteUserRequest)

		user, err := userSvc.DeleteUser(ctx, req.ID)
		return user, err
	}
}
