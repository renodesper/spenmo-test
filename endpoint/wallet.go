package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"gitlab.com/renodesper/spenmo-test/service"
)

type (
	CreateWalletRequest struct {
		TeamID uuid.UUID `json:"teamId"`
		UserID uuid.UUID `json:"userId"`
	}

	GetAllWalletsRequest struct {
		SortBy string `json:"sortBy"`
		Sort   string `json:"sort"`
		Skip   int    `json:"skip"`
		Limit  int    `json:"limit"`
	}

	GetWalletRequest struct {
		ID uuid.UUID `json:"id" validate:"required"`
	}

	UpdateWalletRequest struct {
		ID              uuid.UUID `json:"id" validate:"required"`
		TeamID          uuid.UUID `json:"teamId"`
		UserID          uuid.UUID `json:"userId"`
		BalanceIncrease uint64    `json:"balanceIncrease"`
		BalanceDecrease uint64    `json:"balanceDecrease"`
		DailyLimit      uint64    `json:"dailyLimit"`
		MonthlyLimit    uint64    `json:"monthlyLimit"`
	}

	DeleteWalletByIDRequest struct {
		ID uuid.UUID `json:"id" validate:"required"`
	}

	DeleteWalletsByTeamIDRequest struct {
		TeamID uuid.UUID `json:"teamId" validate:"required"`
	}

	DeleteWalletsByUserIDRequest struct {
		UserID uuid.UUID `json:"userId" validate:"required"`
	}
)

func MakeCreateWalletEndpoint(walletSvc service.WalletService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateWalletRequest)

		walletReq := service.CreateWalletRequest{
			TeamID: req.TeamID,
			UserID: req.UserID,
		}

		wallet, err := walletSvc.CreateWallet(ctx, &walletReq)
		return wallet, err
	}
}

func MakeGetAllWalletsEndpoint(walletSvc service.WalletService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetAllWalletsRequest)

		wallets, err := walletSvc.GetAllWallets(ctx, req.SortBy, req.Sort, req.Skip, req.Limit)
		if err != nil {
			return nil, err
		}

		response = map[string]interface{}{
			"wallets": wallets,
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

func MakeGetWalletEndpoint(walletSvc service.WalletService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetWalletRequest)

		wallet, err := walletSvc.GetWallet(ctx, req.ID)
		return wallet, err
	}
}

func MakeUpdateWalletEndpoint(walletSvc service.WalletService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateWalletRequest)

		walletPayload := service.UpdateWalletRequest{
			BalanceIncrease: req.BalanceIncrease,
			BalanceDecrease: req.BalanceDecrease,
			DailyLimit:      req.DailyLimit,
			MonthlyLimit:    req.MonthlyLimit,
			TeamID:          req.TeamID,
			UserID:          req.UserID,
		}

		wallet, err := walletSvc.UpdateWallet(ctx, req.ID, &walletPayload)
		return wallet, err
	}
}

func MakeDeleteWalletByIDEndpoint(walletSvc service.WalletService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteWalletByIDRequest)

		wallet, err := walletSvc.DeleteWalletByID(ctx, req.ID)
		return wallet, err
	}
}

func MakeDeleteWalletsByTeamIDEndpoint(walletSvc service.WalletService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteWalletsByTeamIDRequest)

		wallet, err := walletSvc.DeleteWalletsByTeamID(ctx, req.TeamID)
		return wallet, err
	}
}

func MakeDeleteWalletsByUserIDEndpoint(walletSvc service.WalletService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteWalletsByUserIDRequest)

		wallet, err := walletSvc.DeleteWalletsByUserID(ctx, req.UserID)
		return wallet, err
	}
}
