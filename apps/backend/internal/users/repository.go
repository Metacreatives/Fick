package users

import (
	"context"
	"errors"
	database "fick/backend/internal/database/generated"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrUsernameUnavailable = errors.New(
	"username unavailable",
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

func (r *Repository) FindByID(
	ctx context.Context,
	id string,
) (database.GetUserByIDRow, bool, error) {
	databaseID, err := strconv.ParseInt(
		id,
		10,
		64,
	)

	if err != nil {
		return database.GetUserByIDRow{},
			false,
			nil
	}

	row, err := r.queries.GetUserByID(
		ctx,
		databaseID,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return database.GetUserByIDRow{},
			false,
			nil
	}

	if err != nil {
		return database.GetUserByIDRow{},
			false,
			err
	}

	return row, true, nil
}

func (r *Repository) FindByUsernameForAuth(
	ctx context.Context,
	username string,
) (database.User, bool, error) {
	user, err := r.queries.GetUserByUsernameForAuth(
		ctx,
		username,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return database.User{},
			false,
			nil
	}

	if err != nil {
		return database.User{},
			false,
			err
	}

	return user, true, nil
}

func (r *Repository) Create(
	ctx context.Context,
	username string,
	passwordHash string,
) (database.CreateUserRow, error) {
	user, err := r.queries.CreateUser(
		ctx,
		database.CreateUserParams{
			Username:     username,
			PasswordHash: passwordHash,
		},
	)

	if err == nil {
		return user, nil
	}

	var postgresError *pgconn.PgError

	if errors.As(err, &postgresError) &&
		postgresError.Code == "23505" {
		return database.CreateUserRow{},
			ErrUsernameUnavailable
	}

	return database.CreateUserRow{}, err
}
