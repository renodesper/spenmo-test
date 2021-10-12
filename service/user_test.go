package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAllUsersIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cases := map[string]struct {
		SortBy        string
		Sort          string
		Skip          int
		Limit         int
		ExpectedError bool
	}{
		"SuccessGetAllUsers": {
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
			svc := NewUserService(Log, DB)

			users, err := svc.GetAllUsers(ctx, tc.SortBy, tc.Sort, tc.Skip, tc.Limit)

			if tc.ExpectedError {
				assert.Error(t, err)
				assert.Empty(t, users)
			} else {
				assert.NoError(t, err)

				if v == "FailedEmptyResponse" {
					assert.NotNil(t, users)
				} else {
					assert.NotEmpty(t, users)
				}
			}
		})
	}
}

func TestGetUserIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cases := map[string]struct {
		ID            string
		ExpectedError bool
	}{
		"SuccessGetUser": {
			ID:            "8e159833-5078-4b0a-80a0-363d82bafd60",
			ExpectedError: false,
		},
		"FailedGetUserNotFound": {
			ID:            "8e159833-5078-4b0a-80a0-363d82bafd61",
			ExpectedError: true,
		},
	}

	for v, tc := range cases {
		t.Run(v, func(t *testing.T) {
			ctx := context.Background()
			svc := NewUserService(Log, DB)

			id, err := uuid.Parse(tc.ID)
			assert.NoError(t, err)

			user, err := svc.GetUser(ctx, id)

			if tc.ExpectedError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
			}
		})
	}
}

func TestCreateUserIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	reinitializeDB()

	cases := map[string]struct {
		Email         string
		Name          string
		ExpectedError bool
	}{
		"SuccessCreateUser": {
			Email:         "user100@example.com",
			Name:          "user100",
			ExpectedError: false,
		},
		"FailedUserExist": {
			Email:         "user100@example.com",
			Name:          "user100",
			ExpectedError: true,
		},
	}

	for v, tc := range cases {
		t.Run(v, func(t *testing.T) {
			ctx := context.Background()
			svc := NewUserService(Log, DB)

			payload := CreateUserRequest{
				Email: tc.Email,
				Name:  tc.Name,
			}
			user, err := svc.CreateUser(ctx, &payload)

			if tc.ExpectedError {
				assert.Error(t, err)
				assert.Empty(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, user)

				user, err := svc.GetUser(ctx, user.ID)
				assert.NoError(t, err)
				assert.NotEmpty(t, user)
			}
		})
	}
}

func TestUpdateUserIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	reinitializeDB()

	cases := map[string]struct {
		ID            string
		Email         string
		Name          string
		ExpectedError bool
	}{
		"SuccessUpdateUser": {
			ID:            "8e159833-5078-4b0a-80a0-363d82bafd60",
			Email:         "user100@example.com",
			Name:          "user100",
			ExpectedError: false,
		},
		"FailedUserNotFound": {
			ID:            "8e159833-5078-4b0a-80a0-363d82bafd61",
			Email:         "user100@example.com",
			Name:          "user100",
			ExpectedError: true,
		},
	}

	for v, tc := range cases {
		t.Run(v, func(t *testing.T) {
			ctx := context.Background()
			svc := NewUserService(Log, DB)

			userID, err := uuid.Parse(tc.ID)
			assert.NoError(t, err)

			payload := UpdateUserRequest{
				Email: tc.Email,
				Name:  tc.Name,
			}
			user, err := svc.UpdateUser(ctx, userID, &payload)

			if tc.ExpectedError {
				assert.Error(t, err)
				assert.Empty(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, user)

				user, err := svc.GetUser(ctx, user.ID)
				assert.NoError(t, err)
				assert.NotEmpty(t, user)
				assert.Equal(t, tc.Email, user.Email)
				assert.Equal(t, tc.Name, user.Name)
			}
		})
	}
}

func TestDeleteUserIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	reinitializeDB()

	cases := map[string]struct {
		ID            string
		ExpectedError bool
	}{
		"SuccessDeleteAssignedUser": {
			ID:            "8e159833-5078-4b0a-80a0-363d82bafd60",
			ExpectedError: false,
		},
		"SuccessDeleteUser": {
			ID:            "0e49e11c-660c-43c5-954e-ef9e89b45833",
			ExpectedError: false,
		},
		"FailedUserNotFound": {
			ID:            "8e159833-5078-4b0a-80a0-363d82bafd61",
			ExpectedError: true,
		},
	}

	for v, tc := range cases {
		t.Run(v, func(t *testing.T) {
			ctx := context.Background()
			svc := NewUserService(Log, DB)
			teamMemberSvc := NewTeamMemberService(Log, DB)
			walletSvc := NewWalletService(Log, DB)
			cardSvc := NewCardService(Log, DB)

			userID, err := uuid.Parse(tc.ID)
			assert.NoError(t, err)

			user, err := svc.DeleteUser(ctx, userID)

			if tc.ExpectedError {
				assert.Error(t, err)
				assert.Empty(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, user)

				user, err := svc.GetUser(ctx, user.ID)
				assert.NoError(t, err)

				assert.Equal(t, true, user.IsDeleted)

				if v == "SuccessDeleteAssignedUser" {
					// NOTE: Related team members should be deleted
					teamMembersPayload := GetTeamMembersRequest{
						UserID: user.ID,
					}
					teamMembers, err := teamMemberSvc.GetTeamMembers(ctx, &teamMembersPayload)
					assert.NoError(t, err)
					assert.NotEmpty(t, teamMembers)

					for _, teamMember := range teamMembers {
						assert.Equal(t, true, teamMember.IsDeleted)
					}
				} else {
					// NOTE: Related wallets should be deleted
					walletsPayload := GetWalletsRequest{
						UserID: user.ID,
					}
					wallets, err := walletSvc.GetWallets(ctx, &walletsPayload)
					assert.NoError(t, err)
					assert.NotEmpty(t, wallets)

					for _, wallet := range wallets {
						assert.Equal(t, true, wallet.IsDeleted)

						// NOTE: Related cards should be deleted
						cardsPayload := GetCardsRequest{
							WalletID: wallet.ID,
						}
						cards, err := cardSvc.GetCards(ctx, &cardsPayload)
						assert.NoError(t, err)
						assert.NotEmpty(t, cards)

						for _, card := range cards {
							assert.Equal(t, true, card.IsDeleted)
						}
					}
				}
			}
		})
	}
}
