package repository

import (
	"time"

	"github.com/google/uuid"
)

type (
	Team struct {
		tableName struct{} `pg:"team"` //nolint

		ID        uuid.UUID `db:"id" json:"id"`
		Name      string    `db:"name" json:"name"`
		IsDeleted bool      `db:"isDeleted" json:"isDeleted"`
		CreatedAt time.Time `db:"createdAt" json:"createdAt"`
		UpdatedAt time.Time `db:"updatedAt" json:"updatedAt"`
	}
)

// MarshalBinary ...
func (u *Team) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

// UnmarshalBinary ...
func (u *Team) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
