package repository

import (
	"context"
	modelDB "hora-server/model"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RepoMessage struct {
	db  *sqlx.DB
	rdb *redis.Client
}

func NewMessageRepository(db *sqlx.DB, rdb *redis.Client) *RepoMessage {
	return &RepoMessage{
		db:  db,
		rdb: rdb,
	}
}

func (r *RepoMessage) Create(ctx context.Context, message *modelDB.MessageDB) (*string, error) {
	message.ID = uuid.New()
	now := time.Now()

	query := `
		INSERT INTO messages (id, content, space_id, user_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query, message.ID, message.Content, message.SpaceID, message.UserID, now)
	if err != nil {
		return nil, err
	}

	idStr := message.ID.String()

	return &idStr, nil
}

func (r *RepoMessage) GetMessages(ctx context.Context) ([]*modelDB.MessageDB, error) {
	const query = `
		SELECT id, content, space_id, user_id, created_at
		FROM messages
		ORDER BY created_at DESC
	`

	var messages []*modelDB.MessageDB
	err := r.db.SelectContext(ctx, &messages, query)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *RepoMessage) PublishMessage(ctx context.Context, spaceID string, data []byte) error {
	err := r.rdb.Publish(ctx, spaceID, data).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RepoMessage) SubscribeMessage(ctx context.Context, spaceID string) *redis.PubSub {
	return r.rdb.Subscribe(ctx, spaceID)
}
