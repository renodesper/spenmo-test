package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/renodesper/spenmo-test/repository"
)

func TestGetAllCardsIntegration(t *testing.T) {
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
		"SuccessGetAllCards": {
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
			svc := NewCardService(Log, DB)

			cards, err := svc.GetAllCards(ctx, tc.SortBy, tc.Sort, tc.Skip, tc.Limit)

			if tc.ExpectedError {
				assert.Error(t, err)
				assert.Empty(t, cards)
			} else {
				assert.NoError(t, err)

				if v == "FailedEmptyResponse" {
					assert.NotNil(t, cards)
				} else {
					assert.NotEmpty(t, cards)
				}
			}
		})
	}
}

func TestGetCardsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cases := map[string]struct {
		WalletID      string
		SortBy        string
		Sort          string
		Skip          int
		Limit         int
		ExpectedError bool
	}{
		"SuccessGetCards": {
			WalletID:      "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			SortBy:        "created_at",
			Sort:          "DESC",
			Skip:          0,
			Limit:         10,
			ExpectedError: false,
		},
		"FailedOffsetMustNotBeNegative": {
			WalletID:      "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			SortBy:        "created_at",
			Sort:          "DESC",
			Skip:          -1,
			Limit:         0,
			ExpectedError: true,
		},
		"FailedLimitMustNotBeNegative": {
			WalletID:      "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			SortBy:        "created_at",
			Sort:          "DESC",
			Skip:          0,
			Limit:         -1,
			ExpectedError: true,
		},
		"FailedEmptyResponse": {
			WalletID:      "d4a6607a-1af7-4571-bdff-2672be72ba0e",
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
			svc := NewCardService(Log, DB)

			payload := GetCardsRequest{
				SortBy: tc.SortBy,
				Sort:   tc.Sort,
				Skip:   tc.Skip,
				Limit:  tc.Limit,
			}

			walletID, err := uuid.Parse(tc.WalletID)
			assert.NoError(t, err)
			payload.WalletID = walletID

			cards, err := svc.GetCards(ctx, &payload)

			if tc.ExpectedError {
				assert.Error(t, err)
				assert.Nil(t, cards)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cards)
			}
		})
	}
}

func TestGetCardIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cases := map[string]struct {
		ID            string
		ExpectedError bool
	}{
		"SuccessGetCard": {
			ID:            "cb21fe95-37e7-4e67-aac3-1b633fe1036d",
			ExpectedError: false,
		},
		"FailedGetCardNotFound": {
			ID:            "cb21fe95-37e7-4e67-aac3-1b633fe1036e",
			ExpectedError: true,
		},
	}

	for v, tc := range cases {
		t.Run(v, func(t *testing.T) {
			ctx := context.Background()
			svc := NewCardService(Log, DB)

			id, err := uuid.Parse(tc.ID)
			assert.NoError(t, err)

			card, err := svc.GetCard(ctx, id)

			if tc.ExpectedError {
				assert.Error(t, err)
				assert.Nil(t, card)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, card)
			}
		})
	}
}

func TestCreateCardIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	reinitializeDB()

	cases := map[string]struct {
		WalletID      string
		CardNo        string
		ExpiryMonth   string
		ExpiryYear    string
		CVV           string
		ExpectedError bool
	}{
		"SuccessCreateCard": {
			WalletID:      "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			CardNo:        "4012888888881881",
			ExpiryMonth:   "12",
			ExpiryYear:    "2022",
			CVV:           "123",
			ExpectedError: false,
		},
		"FailedInvalidCardNumber": {
			WalletID:      "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			CardNo:        "4012888888881882",
			ExpiryMonth:   "12",
			ExpiryYear:    "2022",
			CVV:           "123",
			ExpectedError: true,
		},
		"FailedCardExist": {
			WalletID:      "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			CardNo:        "5200828282828210",
			ExpiryMonth:   "12",
			ExpiryYear:    "2022",
			CVV:           "123",
			ExpectedError: true,
		},
		"FailedInvalidExpiry": {
			WalletID:      "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			CardNo:        "5555555555554444",
			ExpiryMonth:   "13",
			ExpiryYear:    "2020",
			CVV:           "123",
			ExpectedError: true,
		},
		"FailedInvalidCVV1": {
			WalletID:      "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			CardNo:        "5555555555554444",
			ExpiryMonth:   "12",
			ExpiryYear:    "2022",
			CVV:           "12",
			ExpectedError: true,
		},
		"FailedInvalidCVV2": {
			WalletID:      "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			CardNo:        "5555555555554444",
			ExpiryMonth:   "12",
			ExpiryYear:    "2022",
			CVV:           "12345",
			ExpectedError: true,
		},
		"FailedInvalidCVV3": {
			WalletID:      "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			CardNo:        "5555555555554444",
			ExpiryMonth:   "12",
			ExpiryYear:    "2022",
			CVV:           "ABC",
			ExpectedError: true,
		},
		"SuccessCreateCard2": {
			WalletID:      "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			CardNo:        "5555555555554444",
			ExpiryMonth:   "9",
			ExpiryYear:    "2022",
			CVV:           "123",
			ExpectedError: false,
		},
	}

	for v, tc := range cases {
		t.Run(v, func(t *testing.T) {
			ctx := context.Background()
			svc := NewCardService(Log, DB)

			payload := CreateCardRequest{
				CardNo:      tc.CardNo,
				ExpiryMonth: tc.ExpiryMonth,
				ExpiryYear:  tc.ExpiryYear,
				CVV:         tc.CVV,
			}

			var err error
			var walletID uuid.UUID

			if tc.WalletID == "" {
				walletID = uuid.Nil
			} else {
				walletID, err = uuid.Parse(tc.WalletID)
				assert.NoError(t, err)
			}

			payload.WalletID = walletID

			card, err := svc.CreateCard(ctx, &payload)

			if tc.ExpectedError {
				assert.Error(t, err)
				assert.Empty(t, card)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, card)

				card, err := svc.GetCard(ctx, card.ID)
				assert.NoError(t, err)
				assert.NotEmpty(t, card)
			}
		})
	}
}

func TestUpdateCardIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	reinitializeDB()

	ctx := context.Background()
	svc := NewCardService(Log, DB)

	cases := map[string]struct {
		ID            string
		WalletID      string
		CardNo        string
		ExpiryMonth   string
		ExpiryYear    string
		CVV           string
		DailyLimit    uint64
		MonthlyLimit  uint64
		ExpectedError bool
	}{
		"SuccessUpdateUser": {
			ID:            "cb21fe95-37e7-4e67-aac3-1b633fe1036d",
			WalletID:      "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			CardNo:        "5105105105105100",
			ExpiryMonth:   "12",
			ExpiryYear:    "2022",
			CVV:           "1234",
			DailyLimit:    200000,
			MonthlyLimit:  300000,
			ExpectedError: false,
		},
		"FailedUnableToFetchCard": {
			ID:            "cb21fe95-37e7-4e67-aac3-1b633fe1036e",
			WalletID:      "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			CardNo:        "5200828282828210",
			ExpiryMonth:   "12",
			ExpiryYear:    "2022",
			CVV:           "1234",
			DailyLimit:    200000,
			MonthlyLimit:  300000,
			ExpectedError: true,
		},
		"FailedInvalidCardOwner": {
			ID:            "f7bee42b-c2b6-4d9b-b8bb-c9dff69240a4",
			WalletID:      "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			CardNo:        "5200828282828210",
			ExpiryMonth:   "12",
			ExpiryYear:    "2022",
			CVV:           "1234",
			DailyLimit:    200000,
			MonthlyLimit:  300000,
			ExpectedError: true,
		},
		"FailedInvalidCardNumber": {
			ID:            "cb21fe95-37e7-4e67-aac3-1b633fe1036d",
			WalletID:      "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			CardNo:        "5200828282828211",
			ExpiryMonth:   "12",
			ExpiryYear:    "2022",
			CVV:           "1234",
			DailyLimit:    200000,
			MonthlyLimit:  300000,
			ExpectedError: true,
		},
		// NOTE: It looks like there is a race condition when running the integration test that makes the below case failed intermittently
		// "FailedCardExist": {
		// 	ID:            "cb21fe95-37e7-4e67-aac3-1b633fe1036d",
		// 	WalletID:      "d4a6607a-1af7-4571-bdff-2672be72ba0e",
		// 	CardNo:        "5200828282828210",
		// 	ExpiryMonth:   "12",
		// 	ExpiryYear:    "2022",
		// 	CVV:           "1234",
		// 	DailyLimit:    200000,
		// 	MonthlyLimit:  300000,
		// 	ExpectedError: true,
		// },
		"FailedInvalidExpiryMonth": {
			ID:            "cb21fe95-37e7-4e67-aac3-1b633fe1036d",
			WalletID:      "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			CardNo:        "378282246310005",
			ExpiryMonth:   "13",
			ExpiryYear:    "2022",
			CVV:           "1234",
			DailyLimit:    200000,
			MonthlyLimit:  300000,
			ExpectedError: true,
		},
		"FailedInvalidExpiryYear": {
			ID:            "cb21fe95-37e7-4e67-aac3-1b633fe1036d",
			WalletID:      "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			CardNo:        "378282246310005",
			ExpiryMonth:   "12",
			ExpiryYear:    "2020",
			CVV:           "1234",
			DailyLimit:    200000,
			MonthlyLimit:  300000,
			ExpectedError: true,
		},
		"FailedInvalidCVV": {
			ID:            "cb21fe95-37e7-4e67-aac3-1b633fe1036d",
			WalletID:      "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			CardNo:        "378282246310005",
			ExpiryMonth:   "12",
			ExpiryYear:    "2022",
			CVV:           "12345",
			DailyLimit:    200000,
			MonthlyLimit:  300000,
			ExpectedError: true,
		},
		"FailedInvalidDailyLimit": {
			ID:            "cb21fe95-37e7-4e67-aac3-1b633fe1036d",
			WalletID:      "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			CardNo:        "378282246310005",
			ExpiryMonth:   "12",
			ExpiryYear:    "2022",
			CVV:           "1234",
			DailyLimit:    5000,
			MonthlyLimit:  300000,
			ExpectedError: true,
		},
		"FailedInvalidMonthlyLimit": {
			ID:            "cb21fe95-37e7-4e67-aac3-1b633fe1036d",
			WalletID:      "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			CardNo:        "378282246310005",
			ExpiryMonth:   "12",
			ExpiryYear:    "2022",
			CVV:           "1234",
			DailyLimit:    300000,
			MonthlyLimit:  5000,
			ExpectedError: true,
		},
	}

	for v, tc := range cases {
		t.Run(v, func(t *testing.T) {
			cardID, err := uuid.Parse(tc.ID)
			assert.NoError(t, err)

			payload := UpdateCardRequest{
				CardNo:       tc.CardNo,
				ExpiryMonth:  tc.ExpiryMonth,
				ExpiryYear:   tc.ExpiryYear,
				CVV:          tc.CVV,
				DailyLimit:   tc.DailyLimit,
				MonthlyLimit: tc.MonthlyLimit,
			}

			var walletID uuid.UUID

			if tc.WalletID == "" {
				walletID = uuid.Nil
			} else {
				walletID, err = uuid.Parse(tc.WalletID)
				assert.NoError(t, err)
			}

			payload.WalletID = walletID

			card, err := svc.UpdateCard(ctx, cardID, &payload)

			if tc.ExpectedError {
				assert.Error(t, err)
				assert.Empty(t, card)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, card)
			}
		})
	}
}

func TestDeleteCardByIDIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	reinitializeDB()

	cases := map[string]struct {
		CardID        string
		ExpectedError bool
	}{
		"SuccessDeleteCardByID": {
			CardID:        "cb21fe95-37e7-4e67-aac3-1b633fe1036d",
			ExpectedError: false,
		},
		"FailedCardNotFound": {
			CardID:        "cb21fe95-37e7-4e67-aac3-1b633fe1036e",
			ExpectedError: true,
		},
	}

	for v, tc := range cases {
		t.Run(v, func(t *testing.T) {
			ctx := context.Background()
			svc := NewCardService(Log, DB)

			cardID, err := uuid.Parse(tc.CardID)
			assert.NoError(t, err)

			card, err := svc.DeleteCardByID(ctx, cardID)

			if v == "FailedCardNotFound" {
				assert.Error(t, err)
				assert.Nil(t, card)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, card)
			}
		})
	}
}

func TestDeleteCardsByWalletIDIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	reinitializeDB()

	cases := map[string]struct {
		CardID        string
		ExpectedError bool
	}{
		"SuccessDeleteCardsByWalletID": {
			CardID:        "d4a6607a-1af7-4571-bdff-2672be72ba0e",
			ExpectedError: false,
		},
		"FailedCardsNotFound": {
			CardID:        "d4a6607a-1af7-4571-bdff-2672be72ba0f",
			ExpectedError: true,
		},
	}

	for v, tc := range cases {
		t.Run(v, func(t *testing.T) {
			ctx := context.Background()
			svc := NewCardService(Log, DB)

			cardID, err := uuid.Parse(tc.CardID)
			assert.NoError(t, err)

			cards, err := svc.DeleteCardsByWalletID(ctx, cardID)

			if v == "FailedCardsNotFound" {
				assert.NoError(t, err)
				assert.Equal(t, []repository.Card(nil), cards)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, cards)

				for _, card := range cards {
					assert.Equal(t, true, card.IsDeleted)
				}
			}
		})
	}
}
