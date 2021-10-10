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

func decodeCreateWalletRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.CreateWalletRequest

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

func decodeGetAllWalletsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	sortBy := "created_at"
	sort := "DESC"
	skip := 0
	limit := 10

	var req endpoint.GetAllWalletsRequest

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

func decodeGetWalletRequest(_ context.Context, r *http.Request) (interface{}, error) {
	IDStr := bone.GetValue(r, "id")

	var req endpoint.GetWalletRequest

	ID, err := uuid.Parse(IDStr)
	if err != nil {
		return nil, err
	}

	req.ID = ID

	return req, nil
}

func decodeUpdateWalletRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.UpdateWalletRequest

	IDStr := bone.GetValue(r, "id")
	ID, err := uuid.Parse(IDStr)
	if err != nil {
		return nil, err
	}

	req.ID = ID

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

func decodeDeleteWalletByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.DeleteWalletByIDRequest

	IDStr := bone.GetValue(r, "id")
	ID, err := uuid.Parse(IDStr)
	if err != nil {
		return nil, err
	}

	req.ID = ID

	return req, nil
}

func decodeDeleteWalletsByTeamIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.DeleteWalletsByTeamIDRequest

	teamIDStr := bone.GetValue(r, "teamId")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		return nil, err
	}

	req.TeamID = teamID

	return req, nil
}

func decodeDeleteWallestByUserIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.DeleteWalletsByUserIDRequest

	userIDStr := bone.GetValue(r, "userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	req.UserID = userID

	return req, nil
}
