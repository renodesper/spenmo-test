package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"gitlab.com/renodesper/spenmo-test/service"
)

type (
	CreateCardRequest struct {
		CardNo      string    `json:"cardNo" validate:"required"`
		ExpiryMonth string    `json:"expiryMonth" validate:"required"`
		ExpiryYear  string    `json:"expiryYear" validate:"required"`
		CVV         string    `json:"cvv" validate:"required"`
		WalletID    uuid.UUID `json:"walletId" validate:"required"`
	}

	GetAllCardsRequest struct {
		SortBy string `json:"sortBy"`
		Sort   string `json:"sort"`
		Skip   int    `json:"skip"`
		Limit  int    `json:"limit"`
	}

	GetCardRequest struct {
		ID uuid.UUID `json:"id" validate:"required"`
	}

	UpdateCardRequest struct {
		ID           uuid.UUID `json:"id" validate:"required"`
		WalletID     uuid.UUID `json:"walletId"`
		CardNo       string    `json:"cardNo"`
		ExpiryMonth  string    `json:"expiryMonth"`
		ExpiryYear   string    `json:"expiryYear"`
		CVV          string    `json:"cvv"`
		DailyLimit   uint64    `json:"dailyLimit"`
		MonthlyLimit uint64    `json:"monthlyLimit"`
	}

	DeleteCardByIDRequest struct {
		ID uuid.UUID `json:"id" validate:"required"`
	}

	DeleteCardsByWalletIDRequest struct {
		WalletID uuid.UUID `json:"walletId" validate:"required"`
	}
)

func MakeCreateCardEndpoint(CardSvc service.CardService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateCardRequest)

		cardReq := service.CreateCardRequest{
			CardNo:      req.CardNo,
			ExpiryMonth: req.ExpiryMonth,
			ExpiryYear:  req.ExpiryYear,
			CVV:         req.CVV,
			WalletID:    req.WalletID,
		}

		card, err := CardSvc.CreateCard(ctx, &cardReq)
		return card, err
	}
}

func MakeGetAllCardsEndpoint(CardSvc service.CardService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetAllCardsRequest)

		cards, err := CardSvc.GetAllCards(ctx, req.SortBy, req.Sort, req.Skip, req.Limit)
		if err != nil {
			return nil, err
		}

		response = map[string]interface{}{
			"cards": cards,
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

func MakeGetCardEndpoint(CardSvc service.CardService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetCardRequest)

		card, err := CardSvc.GetCard(ctx, req.ID)
		return card, err
	}
}

func MakeUpdateCardEndpoint(CardSvc service.CardService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateCardRequest)

		cardPayload := service.UpdateCardRequest{
			WalletID:     req.WalletID,
			CardNo:       req.CardNo,
			ExpiryMonth:  req.ExpiryMonth,
			ExpiryYear:   req.ExpiryYear,
			CVV:          req.CVV,
			DailyLimit:   req.DailyLimit,
			MonthlyLimit: req.MonthlyLimit,
		}

		card, err := CardSvc.UpdateCard(ctx, req.ID, &cardPayload)
		return card, err
	}
}

func MakeDeleteCardByIDEndpoint(CardSvc service.CardService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteCardByIDRequest)

		card, err := CardSvc.DeleteCardByID(ctx, req.ID)
		return card, err
	}
}

func MakeDeleteCardsByWalletIDEndpoint(CardSvc service.CardService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteCardsByWalletIDRequest)

		cards, err := CardSvc.DeleteCardsByWalletID(ctx, req.WalletID)
		return cards, err
	}
}
