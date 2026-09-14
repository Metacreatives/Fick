package users

import (
	"context"
	"errors"
	database "fick/backend/internal/database/generated"
	"unicode/utf8"
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

var ErrInvalidRegistration = errors.New(
	"invalid registration",
)

func registerUser(
	ctx context.Context,
	repository *Repository,
	username string,
	password string,
) (database.CreateUserRow, error) {
	usernameLength := utf8.RuneCountInString(username)
	passwordLength := utf8.RuneCountInString(password)

	if usernameLength < 1 ||
		usernameLength > 64 ||
		passwordLength < 15 ||
		passwordLength > 256 {
		return database.CreateUserRow{},
			ErrInvalidRegistration
	}

	passwordHash, err := hashPassword(password)
	if err != nil {
		return database.CreateUserRow{}, err
	}

	return repository.Create(
		ctx,
		username,
		passwordHash,
	)
}
