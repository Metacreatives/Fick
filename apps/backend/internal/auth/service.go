package auth

import (
	"context"
	"errors"
	database "fick/backend/internal/database/generated"
	"fick/backend/internal/sessions"
	"fick/backend/internal/users"
	"unicode/utf8"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRegistrationFailed = errors.New(
		"registration failed",
	)

	ErrInvalidCredentials = errors.New(
		"invalid credentials",
	)
)

type Service struct {
	pool *pgxpool.Pool

	users    *users.Repository
	sessions *sessions.Repository

	dummyPasswordHash string
}

func NewService(
	pool *pgxpool.Pool,
	queries *database.Queries,
) (*Service, error) {
	dummyPasswordHash, err := hashPassword(
		"fick-dummy-password-for-authentication",
	)

	if err != nil {
		return nil, err
	}

	return &Service{
		pool: pool,

		users: users.NewRepository(
			queries,
		),

		sessions: sessions.NewRepository(
			queries,
		),

		dummyPasswordHash: dummyPasswordHash,
	}, nil
}

type RegistrationResult struct {
	User  database.CreateUserRow
	Token string
}

func (s *Service) Register(
	ctx context.Context,
	username string,
	password string,
) (RegistrationResult, error) {
	if !validRegistration(
		username,
		password,
	) {
		return RegistrationResult{},
			ErrRegistrationFailed
	}

	passwordHash, err := hashPassword(
		password,
	)

	if err != nil {
		return RegistrationResult{}, err
	}

	sessionToken, err := sessions.NewToken()

	if err != nil {
		return RegistrationResult{}, err
	}

	tx, err := s.pool.Begin(ctx)

	if err != nil {
		return RegistrationResult{}, err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	txQueries := database.New(tx)

	userRepository :=
		users.NewRepository(txQueries)

	sessionRepository :=
		sessions.NewRepository(txQueries)

	user, err := userRepository.Create(
		ctx,
		username,
		passwordHash,
	)

	if errors.Is(
		err,
		users.ErrUsernameUnavailable,
	) {
		return RegistrationResult{},
			ErrRegistrationFailed
	}

	if err != nil {
		return RegistrationResult{}, err
	}

	if err := sessionRepository.Create(
		ctx,
		sessionToken.Hash,
		user.ID,
	); err != nil {
		return RegistrationResult{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return RegistrationResult{}, err
	}

	return RegistrationResult{
		User:  user,
		Token: sessionToken.Raw,
	}, nil
}

func (s *Service) Login(
	ctx context.Context,
	username string,
	password string,
) (string, error) {
	user, found, err :=
		s.users.FindByUsernameForAuth(
			ctx,
			username,
		)

	if err != nil {
		return "", err
	}

	passwordHash := s.dummyPasswordHash

	if found {
		passwordHash = user.PasswordHash
	}

	valid, err := verifyPassword(
		password,
		passwordHash,
	)

	if err != nil {
		return "", err
	}

	if !found || !valid {
		return "", ErrInvalidCredentials
	}

	sessionToken, err := sessions.NewToken()

	if err != nil {
		return "", err
	}

	if err := s.sessions.Create(
		ctx,
		sessionToken.Hash,
		user.ID,
	); err != nil {
		return "", err
	}

	return sessionToken.Raw, nil
}

func validRegistration(
	username string,
	password string,
) bool {
	usernameLength :=
		utf8.RuneCountInString(username)

	passwordLength :=
		utf8.RuneCountInString(password)

	return usernameLength >= 1 &&
		usernameLength <= 64 &&
		passwordLength >= 15 &&
		passwordLength <= 256
}
