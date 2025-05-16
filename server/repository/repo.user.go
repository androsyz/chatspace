package repository

import (
	"context"
	"database/sql"
	"fmt"
	"hora-server/model"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RepoUser struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *RepoUser {
	return &RepoUser{
		db: db,
	}
}

func (r *RepoUser) Create(ctx context.Context, user *model.UserDB) error {
	user.ID = uuid.New()
	now := time.Now()

	query := `
		INSERT INTO users (id, username, email, name, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query, user.ID, user.Username, user.Email, user.Name, user.Password, now, now)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *RepoUser) GetByID(ctx context.Context, id string) (*model.UserDB, error) {
	var user model.UserDB
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, sql.ErrNoRows
	}

	return &user, nil
}

func (r *RepoUser) GetByEmail(ctx context.Context, email string) (*model.UserDB, error) {
	var user model.UserDB
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, sql.ErrNoRows
	}

	return &user, nil
}

func (r *RepoUser) Update(ctx context.Context, user *model.UserDB) error {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users
		SET username = $1, email = $2, password = $3, updated_at = $4
		WHERE id = $5
	`

	result, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.Password, user.UpdatedAt, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *RepoUser) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM users WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
