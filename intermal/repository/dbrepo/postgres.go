package dbrepo

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Authenticate a User
func (m *postgresDBRepo) Authenticate(login, testPassword string) (int, int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var id int
	var hashedPassword string
	var access_level int

	row := m.DB.QueryRowContext(ctx, "select id_account, password, access_level from id_accounts where login = $1", login)
	err := row.Scan(&id, &hashedPassword, &access_level)
	if err != nil {
		return id, access_level, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, 0, "", errors.New("Incorrect password")
	} else if err != nil {
		return 0, 0, "", err
	}

	return id, access_level, hashedPassword, nil
}
