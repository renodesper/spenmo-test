package repository

import (
	"time"

	"github.com/google/uuid"
)

type (
	Wallet struct {
		tableName struct{} `pg:"wallet"` //nolint

		ID           uuid.UUID `db:"id" json:"id"`
		Balance      uint64    `db:"balance" json:"balance"`
		DailyLimit   uint64    `db:"dailyLimit" json:"dailyLimit"`
		MonthlyLimit uint64    `db:"monthlyLimit" json:"monthlyLimit"`
		TeamID       uuid.UUID `db:"teamId" json:"teamId"`
		UserID       uuid.UUID `db:"userId" json:"userId"`
		IsDeleted    bool      `db:"isDeleted" json:"isDeleted"`
		CreatedAt    time.Time `db:"createdAt" json:"createdAt"`
		UpdatedAt    time.Time `db:"updatedAt" json:"updatedAt"`
	}
)

// MarshalBinary ...
func (u *Wallet) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

// UnmarshalBinary ...
func (u *Wallet) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
