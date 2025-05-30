package model

import (
	"time"

	"github.com/google/uuid"
)

type SpaceMemberDB struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	SpaceID   uuid.UUID `db:"space_id"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
}
