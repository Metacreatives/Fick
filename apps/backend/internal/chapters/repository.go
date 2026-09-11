package chapters

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

func (r *Repository) FindByWorkAndNumber(
	ctx context.Context,
	workID string,
	chapterNumber string,
) (Chapter, bool, error) {
	databaseWorkID, err := strconv.ParseInt(workID, 10, 64)
	if err != nil {
		return Chapter{}, false, nil
	}

	databaseChapterNumber, err := strconv.ParseInt(chapterNumber, 10, 64)
	if err != nil {
		return Chapter{}, false, nil
	}

	row, err := r.queries.GetChapterByWorkAndNumber(
		ctx,
		database.GetChapterByWorkAndNumberParams{
			WorkID: databaseWorkID,
			Number: databaseChapterNumber,
		},
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return Chapter{}, false, nil
	}

	if err != nil {
		return Chapter{}, false, err
	}

	chapter := Chapter{
		ID:            strconv.FormatInt(row.ID, 10),
		WorkID:        strconv.FormatInt(row.WorkID, 10),
		Number:        strconv.FormatInt(row.Number, 10),
		Title:         row.Title,
		CreatedAt:     row.CreatedAt.Time.Format(time.RFC3339Nano),
		UpdatedAt:     row.UpdatedAt.Time.Format(time.RFC3339Nano),
		ContentFormat: row.ContentFormat,
		ContentRaw:    row.ContentRaw,
	}

	return chapter, true, nil
}
