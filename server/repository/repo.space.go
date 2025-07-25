package repository

import (
	"context"
	"database/sql"
	"errors"
	modelDB "chatspace-server/model"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RepoSpace struct {
	db *sqlx.DB
}

func NewSpaceRepository(db *sqlx.DB) *RepoSpace {
	return &RepoSpace{
		db: db,
	}
}

func (r *RepoSpace) Create(ctx context.Context, space *modelDB.SpaceDB) (*string, error) {
	space.ID = uuid.New()
	now := time.Now()

	query := `
		INSERT INTO spaces (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query, space.ID, space.Name, space.Description, now, now)
	if err != nil {
		return nil, err
	}

	idStr := space.ID.String()

	return &idStr, nil
}

func (r *RepoSpace) GetSpaceByID(ctx context.Context, id string) (*modelDB.SpaceDB, error) {
	const query = `
		SELECT id, name, description, created_at, updated_at
		FROM spaces
		WHERE id = $1
	`

	var space modelDB.SpaceDB
	err := r.db.GetContext(ctx, &space, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &space, nil
}

func (r *RepoSpace) GetSpaces(ctx context.Context) ([]*modelDB.SpaceDB, error) {
	const query = `
		SELECT id, name, description, created_at, updated_at
		FROM spaces
		ORDER BY created_at DESC
	`

	var spaces []*modelDB.SpaceDB
	err := r.db.SelectContext(ctx, &spaces, query)
	if err != nil {
		return nil, err
	}

	return spaces, nil
}

func (r *RepoSpace) CreateSpaceMember(ctx context.Context, spaceMember *modelDB.SpaceMemberDB) error {
	spaceMember.ID = uuid.New()
	now := time.Now()

	query := `
		INSERT INTO space_members (id, user_id, space_id, role, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, spaceMember.ID, spaceMember.UserID, spaceMember.SpaceID, spaceMember.Role, now)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepoSpace) GetSpaceMember(ctx context.Context, spaceID string) ([]*modelDB.SpaceMemberDB, error) {
	const query = `
		SELECT id, user_id, space_id, role, created_at
		FROM space_members
		WHERE space_id = $1
		ORDER BY created_at DESC
	`

	var members []*modelDB.SpaceMemberDB
	err := r.db.SelectContext(ctx, &members, query, spaceID)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (r *RepoSpace) GetMemberBySpaceID(ctx context.Context, spaceID, role string) ([]*modelDB.UserDB, error) {
	const query = `
		SELECT u.id, u.name, u.email
		FROM users u
		LEFT JOIN space_members sm ON u.id = sm.user_id 
		WHERE sm.space_id = $1 and sm.role = $2
		ORDER BY sm.created_at DESC
	`

	var members []*modelDB.UserDB
	err := r.db.SelectContext(ctx, &members, query, spaceID, role)
	if err != nil {
		return nil, err
	}

	return members, nil
}
