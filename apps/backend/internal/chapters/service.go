package chapters

import (
	"context"
	database "fick/backend/internal/database/generated"
)

func findChapterByWorkAndNumber(
	ctx context.Context,
	repository *Repository,
	workID string,
	chapterNumber string,
) (database.Chapter, bool, error) {
	return repository.FindByWorkAndNumber(
		ctx,
		workID,
		chapterNumber,
	)
}
