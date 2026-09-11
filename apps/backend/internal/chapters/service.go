package chapters

import "context"

func findChapterByWorkAndNumber(
	ctx context.Context,
	repository *Repository,
	workID string,
	chapterNumber string,
) (Chapter, bool, error) {
	return repository.FindByWorkAndNumber(
		ctx,
		workID,
		chapterNumber,
	)
}
