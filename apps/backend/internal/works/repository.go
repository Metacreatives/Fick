package works

import (
	"context"
	"errors"
	database "fick/backend/internal/database/generated"
	"strconv"
	"time"

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
) (Work, bool, error) {
	databaseID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return Work{}, false, nil
	}

	row, err := r.queries.GetWorkByID(ctx, databaseID)
	if errors.Is(err, pgx.ErrNoRows) {
		return Work{}, false, nil
	}

	if err != nil {
		return Work{}, false, err
	}

	work := Work{
		ID:           strconv.FormatInt(row.ID, 10),
		Title:        row.Title,
		Summary:      row.Summary,
		CreatedAt:    row.CreatedAt.Time.Format(time.RFC3339Nano),
		UpdatedAt:    row.UpdatedAt.Time.Format(time.RFC3339Nano),
		Visibility:   WorkVisibility(row.Visibility),
		ChapterCount: strconv.FormatInt(row.ChapterCount, 10),
	}

	return work, true, nil
}
