package postgre

import (
	"context"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"gitlab.com/renodesper/spenmo-test/repository"
	"gitlab.com/renodesper/spenmo-test/util/errors"
	"gitlab.com/renodesper/spenmo-test/util/logger"
)

type (
	WalletRepository interface {
		GetAllWallets(ctx context.Context, sortBy string, sort string, skip int, limit int) ([]repository.Wallet, error)
		GetWalletByID(ctx context.Context, walletID uuid.UUID) (*repository.Wallet, error)
		CreateWallet(ctx context.Context, walletPayload *repository.Wallet) (*repository.Wallet, error)
		UpdateWallet(ctx context.Context, walletID uuid.UUID, walletPayload map[string]interface{}) (*repository.Wallet, error)
		DeleteWalletByID(ctx context.Context, walletID uuid.UUID) (*repository.Wallet, error)
		DeleteWalletsByUserID(ctx context.Context, userID uuid.UUID) ([]repository.Wallet, error)
		DeleteWalletsByTeamID(ctx context.Context, teamID uuid.UUID) ([]repository.Wallet, error)
	}

	WalletRepo struct {
		Log logger.Logger
		Db  *pg.DB
	}
)

var (
	walletTable = "wallet"
)

// CreateWalletRepository creates wallet repository
func CreateWalletRepository(log logger.Logger, db *pg.DB) WalletRepository {
	return &WalletRepo{
		Log: log,
		Db:  db,
	}
}

// GetAllWallets ...
func (r *WalletRepo) GetAllWallets(ctx context.Context, sortBy string, sort string, skip int, limit int) ([]repository.Wallet, error) {
	var wallets []repository.Wallet

	if sortBy == "" {
		sortBy = "created_at"
	}
	if sort == "" {
		sort = "DESC"
	}
	order := fmt.Sprintf("%s %s", sortBy, sort)

	err := r.Db.WithContext(ctx).Model(&wallets).Limit(limit).Offset(skip).Order(order).Select()
	if err != nil {
		return nil, errors.FailedWalletsFetch.AppendError(err)
	}

	return wallets, nil
}

// GetWalletByID ...
func (r *WalletRepo) GetWalletByID(ctx context.Context, walletID uuid.UUID) (*repository.Wallet, error) {
	wallet := repository.Wallet{}

	sql := r.Db.WithContext(ctx).Model(&wallet).Where("id = ?", walletID)

	err := sql.Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, errors.FailedNoRows.AppendError(err)
		}

		return nil, errors.FailedWalletFetch.AppendError(err)
	}

	return &wallet, nil
}

// CreateWallet ...
func (r *WalletRepo) CreateWallet(ctx context.Context, walletPayload *repository.Wallet) (*repository.Wallet, error) {
	var wallet repository.Wallet

	_, err := r.Db.WithContext(ctx).Model(walletPayload).Returning("*").Insert(&wallet)
	if err != nil {
		return nil, errors.FailedWalletCreate.AppendError(err)
	}

	return &wallet, nil
}

func (r *WalletRepo) UpdateWallet(ctx context.Context, walletID uuid.UUID, walletPayload map[string]interface{}) (*repository.Wallet, error) {
	walletPayload["updated_at"] = time.Now()

	var wallet repository.Wallet
	_, err := r.Db.WithContext(ctx).Model(&walletPayload).Table(walletTable).Where("id = ?", walletID).Returning("*").Update(&wallet)
	if err != nil {
		return nil, errors.FailedWalletUpdate.AppendError(err)
	}

	return &wallet, nil
}

func (r *WalletRepo) DeleteWalletByID(ctx context.Context, walletID uuid.UUID) (*repository.Wallet, error) {
	walletPayload := map[string]interface{}{
		"is_deleted": true,
		"updated_at": time.Now(),
	}

	var wallet repository.Wallet
	_, err := r.Db.WithContext(ctx).Model(&walletPayload).Table(walletTable).Where("id = ?", walletID).Returning("*").Update(&wallet)
	if err != nil {
		return nil, errors.FailedWalletDelete.AppendError(err)
	}

	return &wallet, nil
}

func (r *WalletRepo) DeleteWalletsByUserID(ctx context.Context, userID uuid.UUID) ([]repository.Wallet, error) {
	walletPayload := map[string]interface{}{
		"is_deleted": true,
		"updated_at": time.Now(),
	}

	var wallets []repository.Wallet
	_, err := r.Db.WithContext(ctx).Model(&walletPayload).Table(walletTable).Where("user_id = ?", userID).Returning("*").Update(&wallets)
	if err != nil {
		return nil, errors.FailedWalletDelete.AppendError(err)
	}

	return wallets, nil
}

func (r *WalletRepo) DeleteWalletsByTeamID(ctx context.Context, teamID uuid.UUID) ([]repository.Wallet, error) {
	walletPayload := map[string]interface{}{
		"is_deleted": true,
		"updated_at": time.Now(),
	}

	var wallets []repository.Wallet
	_, err := r.Db.WithContext(ctx).Model(&walletPayload).Table(walletTable).Where("team_id = ?", teamID).Returning("*").Update(&wallets)
	if err != nil {
		return nil, errors.FailedWalletDelete.AppendError(err)
	}

	return wallets, nil
}
