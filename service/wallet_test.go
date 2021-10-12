package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAllWalletsIntegration(t *testing.T) {
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
		"SuccessGetAllWallets": {
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
			svc := NewWalletService(Log, DB)

			wallets, err := svc.GetAllWallets(ctx, tc.SortBy, tc.Sort, tc.Skip, tc.Limit)

			if tc.ExpectedError {
				assert.Error(t, err)
				assert.Empty(t, wallets)
			} else {
				assert.NoError(t, err)

				if v == "FailedEmptyResponse" {
					assert.NotNil(t, wallets)
				} else {
					assert.NotEmpty(t, wallets)
				}
			}
		})
	}
}

func TestGetWalletsIntegration(t *testing.T) {
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
		"SuccessGetWalletByTeamID": {
			TeamID:        "933efe12-2219-42df-bd51-a2e84888432d",
			SortBy:        "created_at",
			Sort:          "DESC",
			Skip:          0,
			Limit:         10,
			ExpectedError: false,
		},
		"SuccessGetWalletByUserID": {
			UserID:        "0e49e11c-660c-43c5-954e-ef9e89b45833",
			SortBy:        "created_at",
			Sort:          "DESC",
			Skip:          0,
			Limit:         10,
			ExpectedError: false,
		},
		"FailedOffsetMustNotBeNegative": {
			TeamID:        "933efe12-2219-42df-bd51-a2e84888432d",
			SortBy:        "created_at",
			Sort:          "DESC",
			Skip:          -1,
			Limit:         0,
			ExpectedError: true,
		},
		"FailedLimitMustNotBeNegative": {
			TeamID:        "933efe12-2219-42df-bd51-a2e84888432d",
			SortBy:        "created_at",
			Sort:          "DESC",
			Skip:          0,
			Limit:         -1,
			ExpectedError: true,
		},
		"FailedEmptyResponse": {
			TeamID:        "933efe12-2219-42df-bd51-a2e84888432d",
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
			svc := NewWalletService(Log, DB)

			payload := GetWalletsRequest{
				SortBy: tc.SortBy,
				Sort:   tc.Sort,
				Skip:   tc.Skip,
				Limit:  tc.Limit,
			}

			if tc.TeamID != "" {
				teamID, err := uuid.Parse(tc.TeamID)
				assert.NoError(t, err)
				payload.TeamID = teamID
			} else if tc.UserID != "" {
				userID, err := uuid.Parse(tc.UserID)
				assert.NoError(t, err)
				payload.UserID = userID
			}

			wallets, err := svc.GetWallets(ctx, &payload)

			if tc.ExpectedError {
				assert.Error(t, err)
				assert.Nil(t, wallets)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, wallets)
			}
		})
	}
}

func TestGetWalletIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cases := map[string]struct {
		ID            string
		ExpectedError bool
	}{
		"SuccessGetWallet": {
			ID:            "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			ExpectedError: false,
		},
		"FailedGetWalletNotFound": {
			ID:            "d4a6607a-1af7-4571-bdff-2672be72ba0f",
			ExpectedError: true,
		},
	}

	for v, tc := range cases {
		t.Run(v, func(t *testing.T) {
			ctx := context.Background()
			svc := NewWalletService(Log, DB)

			id, err := uuid.Parse(tc.ID)
			assert.NoError(t, err)

			wallet, err := svc.GetWallet(ctx, id)

			if tc.ExpectedError {
				assert.Error(t, err)
				assert.Nil(t, wallet)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, wallet)
			}
		})
	}
}

func TestCreateWalletIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	reinitializeDB()

	cases := map[string]struct {
		TeamID        string
		UserID        string
		ExpectedError bool
	}{
		"SuccessCreateWalletByTeamID": {
			TeamID:        "933efe12-2219-42df-bd51-a2e84888432d",
			ExpectedError: false,
		},
		"SuccessCreateWalletByUserID": {
			UserID:        "0e49e11c-660c-43c5-954e-ef9e89b45833",
			ExpectedError: false,
		},
		"FailedMissingWalletOwner": {
			TeamID:        "",
			UserID:        "",
			ExpectedError: true,
		},
	}

	for v, tc := range cases {
		t.Run(v, func(t *testing.T) {
			ctx := context.Background()
			svc := NewWalletService(Log, DB)

			payload := CreateWalletRequest{}

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

			wallet, err := svc.CreateWallet(ctx, &payload)

			if tc.ExpectedError {
				assert.Error(t, err)
				assert.Empty(t, wallet)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, wallet)

				wallet, err := svc.GetWallet(ctx, wallet.ID)
				assert.NoError(t, err)
				assert.NotEmpty(t, wallet)
			}
		})
	}
}

func TestUpdateWalletIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	reinitializeDB()

	ctx := context.Background()
	svc := NewWalletService(Log, DB)

	cases := map[string]struct {
		ID              string
		TeamID          string
		UserID          string
		BalanceIncrease uint64
		BalanceDecrease uint64
		DailyLimit      uint64
		MonthlyLimit    uint64
		ExpectedError   bool
	}{
		"SuccessUpdateUser": {
			ID:              "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			TeamID:          "933efe12-2219-42df-bd51-a2e84888432d",
			UserID:          "",
			BalanceIncrease: 1000,
			DailyLimit:      200000,
			MonthlyLimit:    300000,
			ExpectedError:   false,
		},
		"FailedMissingWalletOwner": {
			ID:              "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			TeamID:          "",
			UserID:          "",
			BalanceIncrease: 1000,
			DailyLimit:      200000,
			MonthlyLimit:    300000,
			ExpectedError:   true,
		},
		"FailedInvalidWalletOwner": {
			ID:              "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			TeamID:          "933efe12-2219-42df-bd51-a2e84888432d",
			UserID:          "0e49e11c-660c-43c5-954e-ef9e89b45833",
			BalanceIncrease: 1000,
			DailyLimit:      200000,
			MonthlyLimit:    300000,
			ExpectedError:   true,
		},
	}

	for v, tc := range cases {
		t.Run(v, func(t *testing.T) {

			walletID, err := uuid.Parse(tc.ID)
			assert.NoError(t, err)

			currentWallet, err := svc.GetWallet(ctx, walletID)
			assert.NoError(t, err)
			assert.NotEmpty(t, currentWallet)

			currentBalance := currentWallet.Balance

			payload := UpdateWalletRequest{
				BalanceIncrease: tc.BalanceIncrease,
				DailyLimit:      tc.DailyLimit,
				MonthlyLimit:    tc.MonthlyLimit,
			}

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

			wallet, err := svc.UpdateWallet(ctx, walletID, &payload)

			if tc.ExpectedError {
				assert.Error(t, err)
				assert.Empty(t, wallet)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, wallet)

				updatedWallet, err := svc.GetWallet(ctx, walletID)
				assert.NoError(t, err)
				assert.NotEmpty(t, wallet)

				assert.Equal(t, tc.BalanceIncrease, updatedWallet.Balance-currentBalance)
				assert.Equal(t, tc.DailyLimit, wallet.DailyLimit)
				assert.Equal(t, tc.MonthlyLimit, wallet.MonthlyLimit)
			}
		})
	}

	cases = map[string]struct {
		ID              string
		TeamID          string
		UserID          string
		BalanceIncrease uint64
		BalanceDecrease uint64
		DailyLimit      uint64
		MonthlyLimit    uint64
		ExpectedError   bool
	}{
		"FailedWalletNotFound": {
			ID:              "d4a6607a-1af7-4571-bdff-2672be72ba0f",
			TeamID:          "933efe12-2219-42df-bd51-a2e84888432d",
			UserID:          "",
			BalanceIncrease: 1000,
			DailyLimit:      200000,
			MonthlyLimit:    300000,
			ExpectedError:   true,
		},
		"FailedInvalidWalletOwner": {
			ID:              "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			TeamID:          "933efe12-2219-42df-bd51-a2e84888432e",
			UserID:          "",
			BalanceIncrease: 1000,
			DailyLimit:      200000,
			MonthlyLimit:    300000,
			ExpectedError:   true,
		},
		"FailedInvalidWalletBalanceOperation": {
			ID:              "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			TeamID:          "933efe12-2219-42df-bd51-a2e84888432d",
			UserID:          "",
			BalanceIncrease: 1000,
			BalanceDecrease: 1000,
			DailyLimit:      200000,
			MonthlyLimit:    300000,
			ExpectedError:   true,
		},
		"FailedNotEnoughBalance": {
			ID:              "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			TeamID:          "933efe12-2219-42df-bd51-a2e84888432d",
			UserID:          "",
			BalanceDecrease: 10000,
			DailyLimit:      200000,
			MonthlyLimit:    300000,
			ExpectedError:   true,
		},
		"FailedInvalidDailyLimit": {
			ID:              "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			TeamID:          "933efe12-2219-42df-bd51-a2e84888432d",
			UserID:          "",
			BalanceDecrease: 1000,
			DailyLimit:      5000,
			MonthlyLimit:    300000,
			ExpectedError:   true,
		},
		"FailedInvalidMonthlyLimit": {
			ID:              "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			TeamID:          "933efe12-2219-42df-bd51-a2e84888432d",
			UserID:          "",
			BalanceDecrease: 1000,
			DailyLimit:      300000,
			MonthlyLimit:    5000,
			ExpectedError:   true,
		},
	}

	for v, tc := range cases {
		t.Run(v, func(t *testing.T) {
			walletID, err := uuid.Parse(tc.ID)
			assert.NoError(t, err)

			payload := UpdateWalletRequest{
				BalanceIncrease: tc.BalanceIncrease,
				BalanceDecrease: tc.BalanceDecrease,
				DailyLimit:      tc.DailyLimit,
				MonthlyLimit:    tc.MonthlyLimit,
			}

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

			wallet, err := svc.UpdateWallet(ctx, walletID, &payload)

			if tc.ExpectedError {
				assert.Error(t, err)
				assert.Empty(t, wallet)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, wallet)
			}
		})
	}
}

func TestDeleteWalletByIDIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	reinitializeDB()

	cases := map[string]struct {
		ID            string
		ExpectedError bool
	}{
		"SuccessDeleteWallet": {
			ID:            "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			ExpectedError: false,
		},
		"FailedWalletNotFound": {
			ID:            "d4a6607a-1af7-4571-bdff-2672be72ba0f",
			ExpectedError: true,
		},
	}

	for v, tc := range cases {
		t.Run(v, func(t *testing.T) {
			ctx := context.Background()
			svc := NewWalletService(Log, DB)
			cardSvc := NewCardService(Log, DB)

			walletID, err := uuid.Parse(tc.ID)
			assert.NoError(t, err)

			wallet, err := svc.DeleteWalletByID(ctx, walletID)

			if tc.ExpectedError {
				assert.Error(t, err)
				assert.Empty(t, wallet)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, wallet)

				wallet, err := svc.GetWallet(ctx, walletID)
				assert.NoError(t, err)

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
		})
	}
}

func TestDeleteWalletsByUserIDIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	reinitializeDB()

	cases := map[string]struct {
		UserID        string
		ExpectedError bool
	}{
		"SuccessDeleteWalletByUserID": {
			UserID:        "0e49e11c-660c-43c5-954e-ef9e89b45833",
			ExpectedError: false,
		},
	}

	for v, tc := range cases {
		t.Run(v, func(t *testing.T) {
			ctx := context.Background()
			svc := NewWalletService(Log, DB)
			cardSvc := NewCardService(Log, DB)

			userID, err := uuid.Parse(tc.UserID)
			assert.NoError(t, err)

			wallets, err := svc.DeleteWalletsByUserID(ctx, userID)
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
		})
	}
}

func TestDeleteWalletsByTeamIDIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	reinitializeDB()

	cases := map[string]struct {
		TeamID        string
		ExpectedError bool
	}{
		"SuccessDeleteWalletByTeamID": {
			TeamID:        "933efe12-2219-42df-bd51-a2e84888432d",
			ExpectedError: false,
		},
	}

	for v, tc := range cases {
		t.Run(v, func(t *testing.T) {
			ctx := context.Background()
			svc := NewWalletService(Log, DB)
			cardSvc := NewCardService(Log, DB)

			teamID, err := uuid.Parse(tc.TeamID)
			assert.NoError(t, err)

			wallets, err := svc.DeleteWalletsByTeamID(ctx, teamID)
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
		})
	}
}
