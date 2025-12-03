package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type UsersModel struct {
	DB *sql.DB
}

type User struct {
	Id         uuid.UUID
	Created_At time.Time
	Updated_At time.Time
	Name       string
}

// adds a new user record
const insertQry = `
INSERT INTO users (name)
VALUES ($1)

RETURNING id, created_at, updated_at;`

func (m UsersModel) Insert(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, insertQry, user.Name).Scan(&user.Id, &user.Created_At, &user.Updated_At)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_name_key"`:
			return ErrDuplicateUserName
		default:
			return err
		}
	}

	return nil
}

// Gets a user record by name
const getQry = `
SELECT id, created_at, updated_at, name
FROM users
WHERE name = $1;`

func (m UsersModel) Get(name string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err := m.DB.QueryRowContext(ctx, getQry, name).Scan(
		&user.Id,
		&user.Created_At,
		&user.Updated_At,
		&user.Name,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

// Lists all users
const listUsersQry = `
SELECT id, created_at, updated_at, name
FROM users;`

func (m UsersModel) ListUsers() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, listUsersQry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.Id,
			&user.Created_At,
			&user.Updated_At,
			&user.Name,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
