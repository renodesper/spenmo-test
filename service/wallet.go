package service

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"gitlab.com/renodesper/spenmo-test/repository"
	"gitlab.com/renodesper/spenmo-test/repository/postgre"
	"gitlab.com/renodesper/spenmo-test/util/errors"
	"gitlab.com/renodesper/spenmo-test/util/logger"
)

type (
	// WalletService ...
	WalletService interface {
		GetAllWallets(ctx context.Context, sortBy string, sort string, skip int, limit int) ([]repository.Wallet, error)
		GetWallets(ctx context.Context, payload *GetWalletsRequest) ([]repository.Wallet, error)
		GetWallet(ctx context.Context, walletID uuid.UUID) (*repository.Wallet, error)
		CreateWallet(ctx context.Context, payload *CreateWalletRequest) (*repository.Wallet, error)
		UpdateWallet(ctx context.Context, walletID uuid.UUID, payload *UpdateWalletRequest) (*repository.Wallet, error)
		DeleteWalletByID(ctx context.Context, walletID uuid.UUID) (*repository.Wallet, error)
		DeleteWalletsByUserID(ctx context.Context, userID uuid.UUID) ([]repository.Wallet, error)
		DeleteWalletsByTeamID(ctx context.Context, teamID uuid.UUID) ([]repository.Wallet, error)
	}

	WalletSvc struct {
		Log     logger.Logger
		Wallet  postgre.WalletRepository
		CardSvc CardService
	}

	GetWalletsRequest struct {
		TeamID uuid.UUID
		UserID uuid.UUID
		SortBy string
		Sort   string
		Skip   int
		Limit  int
	}

	CreateWalletRequest struct {
		TeamID uuid.UUID
		UserID uuid.UUID
	}

	UpdateWalletRequest struct {
		BalanceIncrease uint64 `json:"balanceIncrease"`
		BalanceDecrease uint64 `json:"balanceDecrease"`
		DailyLimit      uint64
		MonthlyLimit    uint64
		TeamID          uuid.UUID
		UserID          uuid.UUID
	}
)

const WALLET_DAILY_LIMIT_DEFAULT = 500000
const WALLET_MONTHLY_LIMIT_DEFAULT = 500000
const WALLET_LIMIT_MINIMUM = 10000

// NewWalletService creates user service
func NewWalletService(log logger.Logger, db *pg.DB) WalletService {
	walletRepo := postgre.CreateWalletRepository(log, db)
	cardSvc := NewCardService(log, db)

	return &WalletSvc{
		Log:     log,
		Wallet:  walletRepo,
		CardSvc: cardSvc,
	}
}

func (s *WalletSvc) GetAllWallets(ctx context.Context, sortBy string, sort string, skip int, limit int) ([]repository.Wallet, error) {
	var wallets []repository.Wallet

	wallets, err := s.Wallet.GetAllWallets(ctx, sortBy, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	if wallets == nil {
		wallets = []repository.Wallet{}
	}

	return wallets, nil
}

func (s *WalletSvc) GetWallets(ctx context.Context, payload *GetWalletsRequest) ([]repository.Wallet, error) {
	wallets, err := s.Wallet.GetWallets(ctx, payload.TeamID, payload.UserID, payload.SortBy, payload.Sort, payload.Skip, payload.Limit)
	return wallets, err
}

func (s *WalletSvc) GetWallet(ctx context.Context, walletID uuid.UUID) (*repository.Wallet, error) {
	wallet, err := s.Wallet.GetWalletByID(ctx, walletID)
	return wallet, err
}

func (s *WalletSvc) CreateWallet(ctx context.Context, payload *CreateWalletRequest) (*repository.Wallet, error) {
	ID := uuid.New()

	walletPayload := repository.Wallet{
		ID:           ID,
		DailyLimit:   WALLET_DAILY_LIMIT_DEFAULT,
		MonthlyLimit: WALLET_MONTHLY_LIMIT_DEFAULT,
	}

	if payload.TeamID != uuid.Nil {
		walletPayload.TeamID = payload.TeamID
	}

	if payload.UserID != uuid.Nil {
		walletPayload.UserID = payload.UserID
	}

	wallet, err := s.Wallet.CreateWallet(ctx, &walletPayload)
	return wallet, err
}

func (s *WalletSvc) UpdateWallet(ctx context.Context, walletID uuid.UUID, payload *UpdateWalletRequest) (*repository.Wallet, error) {
	if payload.TeamID == uuid.Nil && payload.UserID == uuid.Nil {
		return nil, errors.MissingWalletOwner
	}

	if payload.TeamID != uuid.Nil && payload.UserID != uuid.Nil {
		return nil, errors.InvalidWalletOwner
	}

	existingWallet, err := s.Wallet.GetWalletByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	if payload.TeamID != uuid.Nil && payload.TeamID != existingWallet.TeamID || payload.UserID != uuid.Nil && payload.UserID != existingWallet.UserID {
		return nil, errors.InvalidWalletOwner
	}

	var walletPayload = make(map[string]interface{})

	if payload.BalanceIncrease > 0 && payload.BalanceDecrease > 0 {
		return nil, errors.InvalidWalletBalanceOperation
	}

	if payload.BalanceIncrease > 0 {
		walletPayload["balance"] = existingWallet.Balance + payload.BalanceIncrease
	}

	if payload.BalanceDecrease > 0 {
		balance := existingWallet.Balance - payload.BalanceDecrease

		if balance > existingWallet.Balance {
			return nil, errors.NotEnoughBalance
		}

		walletPayload["balance"] = balance
	}

	// NOTE: Only for team
	if payload.TeamID != uuid.Nil {
		if payload.DailyLimit > 0 && payload.DailyLimit < WALLET_LIMIT_MINIMUM {
			return nil, errors.InvalidDailyLimit
		}

		if payload.MonthlyLimit > 0 && payload.MonthlyLimit < WALLET_LIMIT_MINIMUM {
			return nil, errors.InvalidMonthlyLimit
		}

		if payload.DailyLimit >= WALLET_LIMIT_MINIMUM {
			walletPayload["daily_limit"] = payload.DailyLimit
		}

		if payload.MonthlyLimit >= WALLET_LIMIT_MINIMUM {
			walletPayload["monthly_limit"] = payload.MonthlyLimit
		}
	}

	wallet, err := s.Wallet.UpdateWallet(ctx, walletID, walletPayload)
	return wallet, err
}

func (s *WalletSvc) DeleteWalletByID(ctx context.Context, walletID uuid.UUID) (*repository.Wallet, error) {
	wallet, err := s.Wallet.DeleteWalletByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	_, err = s.CardSvc.DeleteCardsByWalletID(ctx, walletID)
	if err != nil {
		s.Log.Warn(err)
	}

	return wallet, nil
}

func (s *WalletSvc) DeleteWalletsByUserID(ctx context.Context, userID uuid.UUID) ([]repository.Wallet, error) {
	wallets, err := s.Wallet.DeleteWalletsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	err = s.DeleteCardsByWalletID(ctx, wallets)
	if err != nil {
		s.Log.Warn(err)
	}

	return wallets, nil
}

func (s *WalletSvc) DeleteWalletsByTeamID(ctx context.Context, teamID uuid.UUID) ([]repository.Wallet, error) {
	wallets, err := s.Wallet.DeleteWalletsByTeamID(ctx, teamID)
	if err != nil {
		return nil, err
	}

	err = s.DeleteCardsByWalletID(ctx, wallets)
	if err != nil {
		s.Log.Warn(err)
	}

	return wallets, nil
}

func (s *WalletSvc) DeleteCardsByWalletID(ctx context.Context, wallets []repository.Wallet) error {
	deleteCardsError := errors.FailedCardsDelete

	for _, wallet := range wallets {
		_, err := s.CardSvc.DeleteCardsByWalletID(ctx, wallet.ID)
		if err != nil {
			_ = deleteCardsError.AppendError(err)
		}
	}

	return deleteCardsError
}
