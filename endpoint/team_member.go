package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"gitlab.com/renodesper/spenmo-test/service"
)

type (
	GetTeamMembersRequest struct {
		TeamID uuid.UUID `json:"teamId"`
		UserID uuid.UUID `json:"userId"`
		SortBy string    `json:"sortBy"`
		Sort   string    `json:"sort"`
		Skip   int       `json:"skip"`
		Limit  int       `json:"limit"`
	}

	CreateTeamMemberRequest struct {
		TeamID uuid.UUID `json:"teamId" validate:"required"`
		UserID uuid.UUID `json:"userId" validate:"required"`
	}

	DeleteTeamMemberRequest struct {
		TeamID uuid.UUID `json:"teamId" validate:"required"`
		UserID uuid.UUID `json:"userId" validate:"required"`
	}

	DeleteTeamMembersByTeamIDRequest struct {
		TeamID uuid.UUID `json:"teamId" validate:"required"`
	}

	DeleteTeamMembersByUserIDRequest struct {
		UserID uuid.UUID `json:"userId" validate:"required"`
	}
)

func MakeGetTeamMembersEndpoint(teamMemberSvc service.TeamMemberService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetTeamMembersRequest)

		teamReq := service.GetTeamMembersRequest{
			TeamID: req.TeamID,
			UserID: req.UserID,
			SortBy: req.SortBy,
			Sort:   req.Sort,
			Skip:   req.Skip,
			Limit:  req.Limit,
		}

		teamMembers, err := teamMemberSvc.GetTeamMembers(ctx, &teamReq)
		if err != nil {
			return nil, err
		}

		response = map[string]interface{}{
			"teamMembers": teamMembers,
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

func MakeCreateTeamMemberEndpoint(teamMemberSvc service.TeamMemberService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateTeamMemberRequest)

		teamReq := service.CreateTeamMemberRequest{
			TeamID: req.TeamID,
			UserID: req.UserID,
		}

		teamMember, err := teamMemberSvc.CreateTeamMember(ctx, &teamReq)
		return teamMember, err
	}
}

func MakeDeleteTeamMemberEndpoint(teamMemberSvc service.TeamMemberService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteTeamMemberRequest)

		teamMember, err := teamMemberSvc.DeleteTeamMember(ctx, req.TeamID, req.UserID)
		return teamMember, err
	}
}

func MakeDeleteTeamMembersByTeamIDEndpoint(teamMemberSvc service.TeamMemberService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteTeamMembersByTeamIDRequest)

		teamMembers, err := teamMemberSvc.DeleteTeamMembersByTeamID(ctx, req.TeamID)
		return teamMembers, err
	}
}

func MakeDeleteTeamMembersByUserIDEndpoint(teamMemberSvc service.TeamMemberService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteTeamMembersByUserIDRequest)

		teamMembers, err := teamMemberSvc.DeleteTeamMembersByUserID(ctx, req.UserID)
		return teamMembers, err
	}
}
