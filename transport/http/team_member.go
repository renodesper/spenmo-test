package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/go-zoo/bone"
	"github.com/google/uuid"
	"github.com/iancoleman/strcase"
	"gitlab.com/renodesper/spenmo-test/endpoint"
	"gitlab.com/renodesper/spenmo-test/util/errors"
)

func decodeCreateTeamMemberRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.CreateTeamMemberRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.UnparsableJSON
	}
	defer r.Body.Close()

	validate = validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, errors.InvalidRequest
	}

	return req, nil
}

func decodeGetTeamMembersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	sortBy := "created_at"
	sort := "DESC"
	skip := 0
	limit := 10

	var req endpoint.GetTeamMembersRequest

	teamIDStr := r.URL.Query().Get("teamId")
	if teamIDStr != "" {
		teamID, err := uuid.Parse(teamIDStr)
		if err != nil {
			return nil, errors.UnparsableUUID.AppendError(err)
		}
		req.TeamID = teamID
	}

	userIDStr := r.URL.Query().Get("userId")
	if userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, errors.UnparsableUUID.AppendError(err)
		}
		req.UserID = userID
	}

	sortByParam := r.URL.Query().Get("sortBy")
	if sortByParam != "" {
		sortBy = strcase.ToSnake(sortByParam)
	}
	req.SortBy = sortBy

	sortParam := r.URL.Query().Get("sort")
	if sortParam != "" {
		sort = sortParam
	}
	req.Sort = sort

	skipParam := r.URL.Query().Get("skip")
	if skipParam != "" {
		skip, _ = strconv.Atoi(skipParam)
	}
	req.Skip = skip

	limitParam := r.URL.Query().Get("limit")
	if limitParam != "" {
		limit, _ = strconv.Atoi(limitParam)
	}
	req.Limit = limit

	validate = validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, errors.InvalidRequest
	}

	return req, nil
}

func decodeDeleteTeamMemberRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.DeleteTeamMemberRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.UnparsableJSON
	}
	defer r.Body.Close()

	validate = validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, errors.InvalidRequest
	}

	return req, nil
}

func decodeDeleteTeamMembersByTeamIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.DeleteTeamMembersByTeamIDRequest

	teamIDStr := bone.GetValue(r, "teamId")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		return nil, err
	}

	req.TeamID = teamID

	return req, nil
}

func decodeDeleteTeamMembersByUserIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.DeleteTeamMembersByUserIDRequest

	userIDStr := bone.GetValue(r, "userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	req.UserID = userID

	return req, nil
}
