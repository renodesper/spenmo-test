package repository

import (
	"time"

	"github.com/google/uuid"
)

type (
	Card struct {
		tableName struct{} `pg:"card"` //nolint

		ID           uuid.UUID `db:"id" json:"id"`
		CardNo       string    `db:"cardNo" json:"cardNo"`
		ExpiryMonth  string    `db:"expiryMonth" json:"expiryMonth"`
		ExpiryYear   string    `db:"expiryYear" json:"expiryYear"`
		CVV          string    `db:"cvv" json:"cvv"`
		DailyLimit   uint64    `db:"dailyLimit" json:"dailyLimit"`
		MonthlyLimit uint64    `db:"monthlyLimit" json:"monthlyLimit"`
		WalletID     uuid.UUID `db:"walletId" json:"walletId"`
		IsDeleted    bool      `db:"isDeleted" json:"isDeleted"`
		CreatedAt    time.Time `db:"createdAt" json:"createdAt"`
		UpdatedAt    time.Time `db:"updatedAt" json:"updatedAt"`
	}
)

// MarshalBinary ...
func (u *Card) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

// UnmarshalBinary ...
func (u *Card) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
