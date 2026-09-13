package users

import (
	"context"
	database "fick/backend/internal/database/generated"
)

func findUserByID(
	ctx context.Context,
	repository *Repository,
	id string,
) (database.GetUserByIDRow, bool, error) {
	row, found, err := repository.FindByID(
		ctx,
		id,
	)

	if err != nil {
		return database.GetUserByIDRow{}, false, err
	}

	if !found {
		return database.GetUserByIDRow{}, false, nil
	}

	return row, true, nil
}
