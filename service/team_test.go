package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllTeams(t *testing.T) {
	cases := map[string]struct {
		SortBy           string
		Sort             string
		Skip             int
		Limit            int
		ExpectedError    bool
		ExpectedResponse string
	}{
		"SuccessGetAllTeams": {
			SortBy:           "created_at",
			Sort:             "DESC",
			Skip:             0,
			Limit:            10,
			ExpectedError:    false,
			ExpectedResponse: "[{\"id\":\"933efe12-2219-42df-bd51-a2e84888432d\",\"name\":\"team1\",\"isDeleted\":false,\"createdAt\":\"2021-10-10T16:57:23.672096Z\",\"updatedAt\":\"0001-01-01T00:00:00Z\"}]",
		},
		"FailedOffsetMustNotBeNegative": {
			SortBy:           "created_at",
			Sort:             "DESC",
			Skip:             -1,
			Limit:            0,
			ExpectedError:    true,
			ExpectedResponse: "null",
		},
		"FailedLimitMustNotBeNegative": {
			SortBy:           "created_at",
			Sort:             "DESC",
			Skip:             0,
			Limit:            -1,
			ExpectedError:    true,
			ExpectedResponse: "null",
		},
		"FailedEmptyResponse": {
			SortBy:           "created_at",
			Sort:             "DESC",
			Skip:             100000,
			Limit:            0,
			ExpectedError:    false,
			ExpectedResponse: "[]",
		},
	}

	for v, test := range cases {
		t.Run(v, func(t *testing.T) {

			ctx := context.Background()
			svc := NewTeamService(Log, DB)

			teams, err := svc.GetAllTeams(ctx, test.SortBy, test.Sort, test.Skip, test.Limit)

			if test.ExpectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			b, _ := json.Marshal(teams)
			assert.Equal(t, test.ExpectedResponse, string(b))
		})
	}
}
