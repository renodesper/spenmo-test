package repository

import (
	"time"

	"github.com/google/uuid"
)

type (
	User struct {
		tableName struct{} `pg:"user"` //nolint

		ID        uuid.UUID `db:"id" json:"id"`
		Email     string    `db:"email" json:"email" validate:"email"`
		Name      string    `db:"name" json:"name"`
		IsDeleted bool      `db:"isDeleted" json:"isDeleted"`
		CreatedAt time.Time `db:"createdAt" json:"createdAt"`
		UpdatedAt time.Time `db:"updatedAt" json:"updatedAt"`
	}
)

// MarshalBinary ...
func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

// UnmarshalBinary ...
func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
