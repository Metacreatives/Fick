package sessions

import (
	"context"
	database "fick/backend/internal/database/generated"
	"strconv"
)

type Repository struct {
	queries *database.Queries
}

func NewRepository(
	queries *database.Queries,
) *Repository {
	return &Repository{
		queries: queries,
	}
}

func (r *Repository) Create(
	ctx context.Context,
	tokenHash []byte,
	userID int64,
) error {
	return r.queries.CreateSession(
		ctx,
		database.CreateSessionParams{
			TokenHash: tokenHash,
			UserID:    userID,
		},
	)
}

func (r *Repository) Delete(
	ctx context.Context,
	tokenHashes [][]byte,
) error {
	return r.queries.DeleteSessions(
		ctx,
		tokenHashes,
	)
}

func (r *Repository) FindSessionsById(
	ctx context.Context,
	id string,
) ([]database.Session, error) {
	databaseID, err := strconv.ParseInt(
		id,
		10,
		64,
	)

	if err != nil {
		return []database.Session{},
			nil
	}

	return r.queries.GetSessionsByUserID(ctx, databaseID)
}
