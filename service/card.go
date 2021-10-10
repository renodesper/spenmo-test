package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"gitlab.com/renodesper/spenmo-test/repository"
	"gitlab.com/renodesper/spenmo-test/repository/postgre"
	"gitlab.com/renodesper/spenmo-test/util/card"
	e "gitlab.com/renodesper/spenmo-test/util/error"
	"gitlab.com/renodesper/spenmo-test/util/errors"
	"gitlab.com/renodesper/spenmo-test/util/logger"
)

type (
	// CardService ...
	CardService interface {
		GetAllCards(ctx context.Context, sortBy string, sort string, skip int, limit int) ([]repository.Card, error)
		GetCard(ctx context.Context, cardID uuid.UUID) (*repository.Card, error)
		CreateCard(ctx context.Context, payload *CreateCardRequest) (*repository.Card, error)
		UpdateCard(ctx context.Context, cardID uuid.UUID, payload *UpdateCardRequest) (*repository.Card, error)
		DeleteCardByID(ctx context.Context, cardID uuid.UUID) (*repository.Card, error)
		DeleteCardsByWalletID(ctx context.Context, walletID uuid.UUID) ([]repository.Card, error)
	}

	CardSvc struct {
		Log  logger.Logger
		Card postgre.CardRepository
	}

	CreateCardRequest struct {
		CardNo      string
		ExpiryMonth string
		ExpiryYear  string
		CVV         string
		WalletID    uuid.UUID
	}

	UpdateCardRequest struct {
		WalletID     uuid.UUID
		CardNo       string
		ExpiryMonth  string
		ExpiryYear   string
		CVV          string
		DailyLimit   uint64
		MonthlyLimit uint64
	}
)

const CARD_DAILY_LIMIT_DEFAULT = 500000
const CARD_MONTHLY_LIMIT_DEFAULT = 500000
const CARD_LIMIT_MINIMUM = 10000

// NewCardService creates user service
func NewCardService(log logger.Logger, db *pg.DB) CardService {
	cardRepo := postgre.CreateCardRepository(log, db)

	return &CardSvc{
		Log:  log,
		Card: cardRepo,
	}
}

func (s *CardSvc) GetAllCards(ctx context.Context, sortBy string, sort string, skip int, limit int) ([]repository.Card, error) {
	var cards []repository.Card

	cards, err := s.Card.GetAllCards(ctx, sortBy, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	if cards == nil {
		cards = []repository.Card{}
	}

	return cards, nil
}

func (s *CardSvc) GetCard(ctx context.Context, cardID uuid.UUID) (*repository.Card, error) {
	card, err := s.Card.GetCardByID(ctx, cardID)
	return card, err
}

func (s *CardSvc) CreateCard(ctx context.Context, payload *CreateCardRequest) (*repository.Card, error) {
	ID := uuid.New()

	err := card.IsCardNumberValid(payload.CardNo)
	if err != nil {
		return nil, err
	}

	existingCard, err := s.Card.GetCardByCardNo(ctx, payload.CardNo)
	if err != nil {
		if er := err.(e.Error); er.Code != errors.FailedNoRows.Code {
			return nil, err
		}
	}
	if existingCard != nil {
		return nil, errors.FailedCardExist
	}

	err = card.IsExpiryValid(payload.ExpiryMonth, payload.ExpiryYear)
	if err != nil {
		return nil, err
	}

	err = card.IsCVVValid(payload.CVV)
	if err != nil {
		return nil, err
	}

	month, _ := strconv.Atoi(payload.ExpiryMonth)
	if month > 9 {
		payload.ExpiryMonth = fmt.Sprintf("%d", month)
	} else {
		payload.ExpiryMonth = fmt.Sprintf("0%d", month)
	}

	cardPayload := repository.Card{
		ID:           ID,
		CardNo:       payload.CardNo,
		ExpiryMonth:  payload.ExpiryMonth,
		ExpiryYear:   payload.ExpiryYear,
		CVV:          payload.CVV,
		DailyLimit:   CARD_DAILY_LIMIT_DEFAULT,
		MonthlyLimit: CARD_MONTHLY_LIMIT_DEFAULT,
		WalletID:     payload.WalletID,
	}
	card, err := s.Card.CreateCard(ctx, &cardPayload)
	return card, err
}

func (s *CardSvc) UpdateCard(ctx context.Context, cardID uuid.UUID, payload *UpdateCardRequest) (*repository.Card, error) {
	var cardPayload = make(map[string]interface{})

	existingCard, err := s.Card.GetCardByID(ctx, cardID)
	if err != nil {
		return nil, err
	}

	if payload.WalletID != uuid.Nil && payload.WalletID != existingCard.WalletID {
		return nil, errors.InvalidCardOwner
	}

	if payload.CardNo != "" {
		err := card.IsCardNumberValid(payload.CardNo)
		if err != nil {
			return nil, err
		}

		duplicateCard, err := s.Card.GetCardByCardNo(ctx, payload.CardNo)
		if err != nil {
			if er := err.(e.Error); er.Code != errors.FailedNoRows.Code {
				return nil, err
			}
		}
		if duplicateCard != nil {
			return nil, errors.FailedCardExist
		}

		cardPayload["card_no"] = payload.CardNo
	}

	if payload.ExpiryMonth != "" {
		_, err := card.IsExpiryMonthValid(payload.ExpiryMonth)
		if err != nil {
			return nil, err
		}

		cardPayload["expiry_month"] = payload.ExpiryMonth
	}

	if payload.ExpiryYear != "" {
		_, err := card.IsExpiryYearValid(payload.ExpiryYear)
		if err != nil {
			return nil, err
		}

		cardPayload["expiry_year"] = payload.ExpiryYear
	}

	if payload.CVV != "" {
		err := card.IsCVVValid(payload.CVV)
		if err != nil {
			return nil, err
		}

		cardPayload["cvv"] = payload.CVV
	}

	if payload.DailyLimit > 0 && payload.DailyLimit < CARD_LIMIT_MINIMUM {
		return nil, errors.InvalidDailyLimit
	}

	if payload.MonthlyLimit > 0 && payload.MonthlyLimit < CARD_LIMIT_MINIMUM {
		return nil, errors.InvalidMonthlyLimit
	}

	if payload.DailyLimit >= CARD_LIMIT_MINIMUM {
		cardPayload["daily_limit"] = payload.DailyLimit
	}

	if payload.MonthlyLimit >= CARD_LIMIT_MINIMUM {
		cardPayload["monthly_limit"] = payload.MonthlyLimit
	}

	card, err := s.Card.UpdateCard(ctx, cardID, cardPayload)
	return card, err
}

func (s *CardSvc) DeleteCardByID(ctx context.Context, cardID uuid.UUID) (*repository.Card, error) {
	card, err := s.Card.DeleteCardByID(ctx, cardID)
	return card, err
}

func (s *CardSvc) DeleteCardsByWalletID(ctx context.Context, walletID uuid.UUID) ([]repository.Card, error) {
	cards, err := s.Card.DeleteCardsByWalletID(ctx, walletID)
	return cards, err
}
