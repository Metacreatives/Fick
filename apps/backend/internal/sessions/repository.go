package sessions

import (
	"context"
	database "fick/backend/internal/database/generated"
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
