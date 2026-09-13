package users

import (
	"context"
	"errors"
	database "fick/backend/internal/database/generated"
	"strconv"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	queries *database.Queries
}

func NewRepository(queries *database.Queries) *Repository {
	return &Repository{
		queries: queries,
	}
}

func (r *Repository) FindByID(
	ctx context.Context,
	id string,
) (database.GetUserByIDRow, bool, error) {
	databaseID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return database.GetUserByIDRow{}, false, nil
	}

	row, err := r.queries.GetUserByID(ctx, databaseID)
	if errors.Is(err, pgx.ErrNoRows) {
		return database.GetUserByIDRow{}, false, nil
	}

	if err != nil {
		return database.GetUserByIDRow{}, false, err
	}

	return row, true, nil
}
