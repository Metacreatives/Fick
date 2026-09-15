package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argonMemory      uint32 = 19 * 1024
	argonIterations  uint32 = 2
	argonParallelism uint8  = 1

	saltLength = 16
	keyLength  = 32
)

type passwordParameters struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
}

func hashPassword(
	password string,
) (string, error) {
	salt := make([]byte, saltLength)

	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf(
			"generate password salt: %w",
			err,
		)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argonIterations,
		argonMemory,
		argonParallelism,
		keyLength,
	)

	encodedSalt :=
		base64.RawStdEncoding.EncodeToString(salt)

	encodedHash :=
		base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		argonMemory,
		argonIterations,
		argonParallelism,
		encodedSalt,
		encodedHash,
	), nil
}

func verifyPassword(
	password string,
	encoded string,
) (bool, error) {
	parameters, salt, expectedHash, err :=
		parsePasswordHash(encoded)

	if err != nil {
		return false, err
	}

	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		parameters.iterations,
		parameters.memory,
		parameters.parallelism,
		uint32(len(expectedHash)),
	)

	return subtle.ConstantTimeCompare(
		actualHash,
		expectedHash,
	) == 1, nil
}

func parsePasswordHash(
	encoded string,
) (
	passwordParameters,
	[]byte,
	[]byte,
	error,
) {
	parts := strings.Split(encoded, "$")

	if len(parts) != 6 ||
		parts[1] != "argon2id" {
		return passwordParameters{},
			nil,
			nil,
			errors.New(
				"invalid password hash format",
			)
	}

	var version int

	if _, err := fmt.Sscanf(
		parts[2],
		"v=%d",
		&version,
	); err != nil {
		return passwordParameters{},
			nil,
			nil,
			fmt.Errorf(
				"parse password hash version: %w",
				err,
			)
	}

	if version != argon2.Version {
		return passwordParameters{},
			nil,
			nil,
			errors.New(
				"unsupported Argon2 version",
			)
	}

	var memory uint32
	var iterations uint32
	var parallelism uint32

	if _, err := fmt.Sscanf(
		parts[3],
		"m=%d,t=%d,p=%d",
		&memory,
		&iterations,
		&parallelism,
	); err != nil {
		return passwordParameters{},
			nil,
			nil,
			fmt.Errorf(
				"parse password parameters: %w",
				err,
			)
	}

	if memory == 0 ||
		memory > 256*1024 ||
		iterations == 0 ||
		iterations > 10 ||
		parallelism == 0 ||
		parallelism > 16 {
		return passwordParameters{},
			nil,
			nil,
			errors.New(
				"invalid password parameters",
			)
	}

	salt, err :=
		base64.RawStdEncoding.DecodeString(
			parts[4],
		)

	if err != nil {
		return passwordParameters{},
			nil,
			nil,
			fmt.Errorf(
				"decode password salt: %w",
				err,
			)
	}

	expectedHash, err :=
		base64.RawStdEncoding.DecodeString(
			parts[5],
		)

	if err != nil {
		return passwordParameters{},
			nil,
			nil,
			fmt.Errorf(
				"decode password hash: %w",
				err,
			)
	}

	if len(salt) < 8 ||
		len(salt) > 64 ||
		len(expectedHash) < 16 ||
		len(expectedHash) > 64 {
		return passwordParameters{},
			nil,
			nil,
			errors.New(
				"invalid password hash sizes",
			)
	}

	return passwordParameters{
			memory:      memory,
			iterations:  iterations,
			parallelism: uint8(parallelism),
		},
		salt,
		expectedHash,
		nil
}
