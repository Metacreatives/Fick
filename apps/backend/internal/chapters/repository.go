package chapters

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

func (r *Repository) FindByWorkAndNumber(
	ctx context.Context,
	workID string,
	chapterNumber string,
) (database.Chapter, bool, error) {
	databaseWorkID, err := strconv.ParseInt(workID, 10, 64)
	if err != nil {
		return database.Chapter{}, false, nil
	}

	databaseChapterNumber, err := strconv.ParseInt(chapterNumber, 10, 64)
	if err != nil {
		return database.Chapter{}, false, nil
	}

	row, err := r.queries.GetChapterByWorkAndNumber(
		ctx,
		database.GetChapterByWorkAndNumberParams{
			WorkID: databaseWorkID,
			Number: databaseChapterNumber,
		},
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return database.Chapter{}, false, nil
	}

	if err != nil {
		return database.Chapter{}, false, err
	}

	return row, true, nil
}
