package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"gitlab.com/renodesper/spenmo-test/service"
)

type (
	CreateTeamRequest struct {
		Name string `json:"name" validate:"required"`
	}

	GetAllTeamsRequest struct {
		SortBy string `json:"sortBy"`
		Sort   string `json:"sort"`
		Skip   int    `json:"skip"`
		Limit  int    `json:"limit"`
	}

	GetTeamRequest struct {
		ID uuid.UUID `json:"id" validate:"required"`
	}

	UpdateTeamRequest struct {
		ID   uuid.UUID `json:"id" validate:"required"`
		Name string
	}

	DeleteTeamRequest struct {
		ID uuid.UUID `json:"id" validate:"required"`
	}
)

func MakeCreateTeamEndpoint(teamSvc service.TeamService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateTeamRequest)

		teamReq := service.CreateTeamRequest{
			Name: req.Name,
		}

		team, err := teamSvc.CreateTeam(ctx, &teamReq)
		return team, err
	}
}

func MakeGetAllTeamsEndpoint(teamSvc service.TeamService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetAllTeamsRequest)

		teams, err := teamSvc.GetAllTeams(ctx, req.SortBy, req.Sort, req.Skip, req.Limit)
		if err != nil {
			return nil, err
		}

		response = map[string]interface{}{
			"teams": teams,
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

func MakeGetTeamEndpoint(teamSvc service.TeamService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetTeamRequest)

		team, err := teamSvc.GetTeam(ctx, req.ID)
		return team, err
	}
}

func MakeUpdateTeamEndpoint(teamSvc service.TeamService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateTeamRequest)

		teamPayload := service.UpdateTeamRequest{
			Name: req.Name,
		}

		team, err := teamSvc.UpdateTeam(ctx, req.ID, &teamPayload)
		return team, err
	}
}

func MakeDeleteTeamEndpoint(teamSvc service.TeamService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteTeamRequest)

		team, err := teamSvc.DeleteTeam(ctx, req.ID)
		return team, err
	}
}
