package works

import (
	"context"
	responses "fick/backend/internal/api/generated"
	database "fick/backend/internal/database/generated"
)

func findPublicWorkByID(
	ctx context.Context,
	repository *Repository,
	id string,
) (database.GetWorkByIDRow, bool, error) {
	work, found, err := repository.FindByID(
		ctx,
		id,
	)

	if err != nil {
		return database.GetWorkByIDRow{}, false, err
	}

	if !found {
		return database.GetWorkByIDRow{}, false, nil
	}

	workVisibility := responses.WorkVisibility(work.Visibility)

	if workVisibility != responses.Public {
		return database.GetWorkByIDRow{}, false, nil
	}

	return work, true, nil
}
