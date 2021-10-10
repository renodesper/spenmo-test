package repository

import (
	"time"

	"github.com/google/uuid"
)

type (
	TeamMember struct {
		tableName struct{} `pg:"team_member"` //nolint

		ID        uuid.UUID `db:"id" json:"id"`
		TeamID    uuid.UUID `db:"teamId" json:"teamId"`
		UserID    uuid.UUID `db:"userId" json:"userId"`
		IsDeleted bool      `db:"isDeleted" json:"isDeleted"`
		CreatedAt time.Time `db:"createdAt" json:"createdAt"`
		UpdatedAt time.Time `db:"updatedAt" json:"updatedAt"`
	}
)

// MarshalBinary ...
func (u *TeamMember) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

// UnmarshalBinary ...
func (u *TeamMember) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
