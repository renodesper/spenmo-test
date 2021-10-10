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
	CardRepository interface {
		GetAllCards(ctx context.Context, sortBy string, sort string, skip int, limit int) ([]repository.Card, error)
		GetCardByID(ctx context.Context, cardID uuid.UUID) (*repository.Card, error)
		GetCardByCardNo(ctx context.Context, cardNo string) (*repository.Card, error)
		CreateCard(ctx context.Context, cardPayload *repository.Card) (*repository.Card, error)
		UpdateCard(ctx context.Context, cardID uuid.UUID, cardPayload map[string]interface{}) (*repository.Card, error)
		DeleteCardByID(ctx context.Context, cardID uuid.UUID) (*repository.Card, error)
		DeleteCardsByWalletID(ctx context.Context, walletID uuid.UUID) ([]repository.Card, error)
	}

	CardRepo struct {
		Log logger.Logger
		Db  *pg.DB
	}
)

var (
	cardTable = "card"
)

// CreateCardRepository creates card repository
func CreateCardRepository(log logger.Logger, db *pg.DB) CardRepository {
	return &CardRepo{
		Log: log,
		Db:  db,
	}
}

// GetAllCards ...
func (r *CardRepo) GetAllCards(ctx context.Context, sortBy string, sort string, skip int, limit int) ([]repository.Card, error) {
	var cards []repository.Card

	if sortBy == "" {
		sortBy = "created_at"
	}
	if sort == "" {
		sort = "DESC"
	}
	order := fmt.Sprintf("%s %s", sortBy, sort)

	err := r.Db.WithContext(ctx).Model(&cards).Limit(limit).Offset(skip).Order(order).Select()
	if err != nil {
		return nil, errors.FailedCardsFetch.AppendError(err)
	}

	return cards, nil
}

// GetCardByID ...
func (r *CardRepo) GetCardByID(ctx context.Context, cardID uuid.UUID) (*repository.Card, error) {
	card := repository.Card{}

	sql := r.Db.WithContext(ctx).Model(&card).Where("id = ?", cardID)

	err := sql.Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, errors.FailedNoRows.AppendError(err)
		}

		return nil, errors.FailedCardFetch.AppendError(err)
	}

	return &card, nil
}

// GetCardByCardNo ...
func (r *CardRepo) GetCardByCardNo(ctx context.Context, cardNo string) (*repository.Card, error) {
	card := repository.Card{}

	sql := r.Db.WithContext(ctx).Model(&card).Where("card_no = ?", cardNo)

	err := sql.Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, errors.FailedNoRows.AppendError(err)
		}

		return nil, errors.FailedCardFetch.AppendError(err)
	}

	return &card, nil
}

// CreateCard ...
func (r *CardRepo) CreateCard(ctx context.Context, cardPayload *repository.Card) (*repository.Card, error) {
	var card repository.Card

	_, err := r.Db.WithContext(ctx).Model(cardPayload).Returning("*").Insert(&card)
	if err != nil {
		return nil, errors.FailedCardCreate.AppendError(err)
	}

	return &card, nil
}

func (r *CardRepo) UpdateCard(ctx context.Context, cardID uuid.UUID, cardPayload map[string]interface{}) (*repository.Card, error) {
	cardPayload["updated_at"] = time.Now()

	var card repository.Card
	_, err := r.Db.WithContext(ctx).Model(&cardPayload).Table(cardTable).Where("id = ?", cardID).Returning("*").Update(&card)
	if err != nil {
		return nil, errors.FailedCardUpdate.AppendError(err)
	}

	return &card, nil
}

func (r *CardRepo) DeleteCardByID(ctx context.Context, cardID uuid.UUID) (*repository.Card, error) {
	cardPayload := map[string]interface{}{
		"is_deleted": true,
		"updated_at": time.Now(),
	}

	var card repository.Card
	_, err := r.Db.WithContext(ctx).Model(&cardPayload).Table(cardTable).Where("id = ?", cardID).Returning("*").Update(&card)
	if err != nil {
		return nil, errors.FailedCardDelete.AppendError(err)
	}

	return &card, nil
}

func (r *CardRepo) DeleteCardsByWalletID(ctx context.Context, walletID uuid.UUID) ([]repository.Card, error) {
	cardPayload := map[string]interface{}{
		"is_deleted": true,
		"updated_at": time.Now(),
	}

	var cards []repository.Card
	_, err := r.Db.WithContext(ctx).Model(&cardPayload).Table(cardTable).Where("wallet_id = ?", walletID).Returning("*").Update(&cards)
	if err != nil {
		return nil, errors.FailedCardDelete.AppendError(err)
	}

	return cards, nil
}
