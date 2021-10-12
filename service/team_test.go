package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAllTeams(t *testing.T) {
	cases := map[string]struct {
		SortBy        string
		Sort          string
		Skip          int
		Limit         int
		ExpectedError bool
	}{
		"SuccessGetAllTeams": {
			SortBy:        "created_at",
			Sort:          "DESC",
			Skip:          0,
			Limit:         10,
			ExpectedError: false,
		},
		"FailedOffsetMustNotBeNegative": {
			SortBy:        "created_at",
			Sort:          "DESC",
			Skip:          -1,
			Limit:         0,
			ExpectedError: true,
		},
		"FailedLimitMustNotBeNegative": {
			SortBy:        "created_at",
			Sort:          "DESC",
			Skip:          0,
			Limit:         -1,
			ExpectedError: true,
		},
		"FailedEmptyResponse": {
			SortBy:        "created_at",
			Sort:          "DESC",
			Skip:          100000,
			Limit:         0,
			ExpectedError: false,
		},
	}

	for v, tc := range cases {
		t.Run(v, func(t *testing.T) {
			ctx := context.Background()
			svc := NewTeamService(Log, DB)

			teams, err := svc.GetAllTeams(ctx, tc.SortBy, tc.Sort, tc.Skip, tc.Limit)

			if tc.ExpectedError {
				assert.Error(t, err)
				assert.Empty(t, teams)
			} else {
				assert.NoError(t, err)

				if v == "FailedEmptyResponse" {
					assert.NotNil(t, teams)
				} else {
					assert.NotEmpty(t, teams)
				}
			}
		})
	}
}

func TestGetTeam(t *testing.T) {
	cases := map[string]struct {
		ID            string
		ExpectedError bool
	}{
		"SuccessGetTeam": {
			ID:            "933efe12-2219-42df-bd51-a2e84888432d",
			ExpectedError: false,
		},
		"FailedGetTeamNotFound": {
			ID:            "933efe12-2219-42df-bd51-a2e84888432e",
			ExpectedError: true,
		},
	}

	for v, tc := range cases {
		t.Run(v, func(t *testing.T) {
			ctx := context.Background()
			svc := NewTeamService(Log, DB)

			id, err := uuid.Parse(tc.ID)
			assert.NoError(t, err)

			team, err := svc.GetTeam(ctx, id)

			if tc.ExpectedError {
				assert.Error(t, err)
				assert.Nil(t, team)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, team)
			}
		})
	}
}

func TestCreateTeam(t *testing.T) {
	reinitializeDB()

	cases := map[string]struct {
		Name          string
		ExpectedError bool
	}{
		"SuccessCreateTeam": {
			Name:          "team100",
			ExpectedError: false,
		},
		"FailedTeamExist": {
			Name:          "team100",
			ExpectedError: true,
		},
	}

	for v, tc := range cases {
		t.Run(v, func(t *testing.T) {
			ctx := context.Background()
			svc := NewTeamService(Log, DB)

			payload := CreateTeamRequest{
				Name: tc.Name,
			}
			team, err := svc.CreateTeam(ctx, &payload)

			if tc.ExpectedError {
				assert.Error(t, err)
				assert.Empty(t, team)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, team)
			}
		})
	}
}

func TestUpdateTeam(t *testing.T) {
	reinitializeDB()

	cases := map[string]struct {
		ID            string
		Name          string
		ExpectedError bool
	}{
		"SuccessUpdateTeam": {
			ID:            "933efe12-2219-42df-bd51-a2e84888432d",
			Name:          "team100",
			ExpectedError: false,
		},
		"FailedTeamNotFound": {
			ID:            "dd95d24e-c934-4b67-8621-c7637be1dfd0",
			Name:          "team100",
			ExpectedError: true,
		},
	}

	for v, tc := range cases {
		t.Run(v, func(t *testing.T) {
			ctx := context.Background()
			svc := NewTeamService(Log, DB)

			teamID, err := uuid.Parse(tc.ID)
			assert.NoError(t, err)

			payload := UpdateTeamRequest{
				Name: tc.Name,
			}
			team, err := svc.UpdateTeam(ctx, teamID, &payload)

			if tc.ExpectedError {
				assert.Error(t, err)
				assert.Empty(t, team)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, team)
			}
		})
	}
}
