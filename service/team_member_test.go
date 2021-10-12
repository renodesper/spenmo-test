package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetTeamMembersIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cases := map[string]struct {
		TeamID        string
		UserID        string
		SortBy        string
		Sort          string
		Skip          int
		Limit         int
		ExpectedError bool
	}{
		"SuccessGetTeamMember": {
			TeamID:        "933efe12-2219-42df-bd51-a2e84888432d",
			UserID:        "8e159833-5078-4b0a-80a0-363d82bafd60",
			SortBy:        "created_at",
			Sort:          "DESC",
			Skip:          0,
			Limit:         10,
			ExpectedError: false,
		},
		"FailedOffsetMustNotBeNegative": {
			TeamID:        "933efe12-2219-42df-bd51-a2e84888432d",
			UserID:        "8e159833-5078-4b0a-80a0-363d82bafd60",
			SortBy:        "created_at",
			Sort:          "DESC",
			Skip:          -1,
			Limit:         0,
			ExpectedError: true,
		},
		"FailedLimitMustNotBeNegative": {
			TeamID:        "933efe12-2219-42df-bd51-a2e84888432d",
			UserID:        "8e159833-5078-4b0a-80a0-363d82bafd60",
			SortBy:        "created_at",
			Sort:          "DESC",
			Skip:          0,
			Limit:         -1,
			ExpectedError: true,
		},
		"FailedEmptyResponse": {
			TeamID:        "933efe12-2219-42df-bd51-a2e84888432d",
			UserID:        "8e159833-5078-4b0a-80a0-363d82bafd60",
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
			svc := NewTeamMemberService(Log, DB)

			payload := GetTeamMembersRequest{
				SortBy: tc.SortBy,
				Sort:   tc.Sort,
				Skip:   tc.Skip,
				Limit:  tc.Limit,
			}

			teamID, err := uuid.Parse(tc.TeamID)
			assert.NoError(t, err)
			payload.TeamID = teamID

			userID, err := uuid.Parse(tc.UserID)
			assert.NoError(t, err)
			payload.UserID = userID

			teamMembers, err := svc.GetTeamMembers(ctx, &payload)

			if tc.ExpectedError {
				assert.Error(t, err)
				assert.Nil(t, teamMembers)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, teamMembers)
			}
		})
	}
}

func TestCreateTeamMemberIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	reinitializeDB()

	cases := map[string]struct {
		TeamID        string
		UserID        string
		ExpectedError bool
	}{
		"SuccessCreateTeamMember": {
			TeamID:        "933efe12-2219-42df-bd51-a2e84888432d",
			UserID:        "bafdd20b-a2a5-41ca-b4a2-26fccbd029dd",
			ExpectedError: false,
		},
		"FailedMissingTeamID": {
			TeamID:        "",
			UserID:        "bafdd20b-a2a5-41ca-b4a2-26fccbd029dd",
			ExpectedError: true,
		},
		"FailedMissingUserID": {
			TeamID:        "933efe12-2219-42df-bd51-a2e84888432d",
			UserID:        "",
			ExpectedError: true,
		},
		"FailedTeamMemberExist": {
			TeamID:        "933efe12-2219-42df-bd51-a2e84888432d",
			UserID:        "8e159833-5078-4b0a-80a0-363d82bafd60",
			ExpectedError: true,
		},
	}

	for v, tc := range cases {
		t.Run(v, func(t *testing.T) {
			ctx := context.Background()
			svc := NewTeamMemberService(Log, DB)

			payload := CreateTeamMemberRequest{}

			var err error
			var teamID uuid.UUID
			var userID uuid.UUID

			if tc.TeamID == "" {
				teamID = uuid.Nil
			} else {
				teamID, err = uuid.Parse(tc.TeamID)
				assert.NoError(t, err)
			}

			payload.TeamID = teamID

			if tc.UserID == "" {
				userID = uuid.Nil
			} else {
				userID, err = uuid.Parse(tc.UserID)
				assert.NoError(t, err)
			}

			payload.UserID = userID

			teamMember, err := svc.CreateTeamMember(ctx, &payload)

			if tc.ExpectedError {
				assert.Error(t, err)
				assert.Nil(t, teamMember)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, teamMember)
			}
		})
	}
}

func TestDeleteTeamMemberIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	reinitializeDB()

	cases := map[string]struct {
		TeamID        string
		UserID        string
		ExpectedError bool
	}{
		"SuccessDeleteTeamMember": {
			TeamID:        "933efe12-2219-42df-bd51-a2e84888432d",
			UserID:        "8e159833-5078-4b0a-80a0-363d82bafd60",
			ExpectedError: false,
		},
		"FailedTeamMemberNotFound1": {
			TeamID:        "933efe12-2219-42df-bd51-a2e84888432e",
			UserID:        "8e159833-5078-4b0a-80a0-363d82bafd61",
			ExpectedError: true,
		},
		"FailedTeamMemberNotFound2": {
			TeamID:        "",
			UserID:        "8e159833-5078-4b0a-80a0-363d82bafd60",
			ExpectedError: true,
		},
		"FailedTeamMemberNotFound3": {
			TeamID:        "933efe12-2219-42df-bd51-a2e84888432d",
			UserID:        "",
			ExpectedError: true,
		},
	}

	for v, tc := range cases {
		t.Run(v, func(t *testing.T) {
			ctx := context.Background()
			svc := NewTeamMemberService(Log, DB)

			var err error
			var teamID uuid.UUID
			var userID uuid.UUID

			if tc.TeamID == "" {
				teamID = uuid.Nil
			} else {
				teamID, err = uuid.Parse(tc.TeamID)
				assert.NoError(t, err)
			}

			if tc.UserID == "" {
				userID = uuid.Nil
			} else {
				userID, err = uuid.Parse(tc.UserID)
				assert.NoError(t, err)
			}

			teamMember, err := svc.DeleteTeamMember(ctx, teamID, userID)
			if tc.ExpectedError {
				assert.Error(t, err)
				assert.Nil(t, teamMember)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, teamMember)
			}
		})
	}
}

func TestDeleteTeamMembersByTeamIDIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	reinitializeDB()

	cases := map[string]struct {
		TeamID      string
		ExpectedNil bool
	}{
		"SuccessDeleteTeamMember": {
			TeamID:      "933efe12-2219-42df-bd51-a2e84888432d",
			ExpectedNil: false,
		},
		"FailedTeamMemberNotFound1": {
			TeamID:      "933efe12-2219-42df-bd51-a2e84888432e",
			ExpectedNil: true,
		},
		"FailedTeamMemberNotFound2": {
			TeamID:      "",
			ExpectedNil: true,
		},
	}

	for v, tc := range cases {
		t.Run(v, func(t *testing.T) {
			ctx := context.Background()
			svc := NewTeamMemberService(Log, DB)

			var err error
			var teamID uuid.UUID

			if tc.TeamID == "" {
				teamID = uuid.Nil
			} else {
				teamID, err = uuid.Parse(tc.TeamID)
				assert.NoError(t, err)
			}

			teamMember, err := svc.DeleteTeamMembersByTeamID(ctx, teamID)
			if tc.ExpectedNil {
				assert.Nil(t, err)
				assert.Nil(t, teamMember)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, teamMember)
			}
		})
	}
}

func TestDeleteTeamMembersByUserIDIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	reinitializeDB()

	cases := map[string]struct {
		UserID      string
		ExpectedNil bool
	}{
		"SuccessDeleteTeamMember": {
			UserID:      "8e159833-5078-4b0a-80a0-363d82bafd60",
			ExpectedNil: false,
		},
		"FailedTeamMemberNotFound1": {
			UserID:      "8e159833-5078-4b0a-80a0-363d82bafd61",
			ExpectedNil: true,
		},
		"FailedTeamMemberNotFound2": {
			UserID:      "",
			ExpectedNil: true,
		},
	}

	for v, tc := range cases {
		t.Run(v, func(t *testing.T) {
			ctx := context.Background()
			svc := NewTeamMemberService(Log, DB)

			var err error
			var userID uuid.UUID

			if tc.UserID == "" {
				userID = uuid.Nil
			} else {
				userID, err = uuid.Parse(tc.UserID)
				assert.NoError(t, err)
			}

			teamMember, err := svc.DeleteTeamMembersByUserID(ctx, userID)
			if tc.ExpectedNil {
				assert.Nil(t, err)
				assert.Nil(t, teamMember)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, teamMember)
			}
		})
	}
}
