package model

import (
	"time"

	"github.com/google/uuid"
)

type MessageDB struct {
	ID        uuid.UUID `db:"id"`
	Content   string    `db:"content"`
	UserID    uuid.UUID `db:"user_id"`
	SpaceID   uuid.UUID `db:"space_id"`
	CreatedAt time.Time `db:"created_at"`
}
