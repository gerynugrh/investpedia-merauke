package model

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
)

type UserModel struct {
	ID 			int
	username	string
	password	string
}

type UserRepository struct {
	Conn *sql.DB
}

func (m *UserRepository) Insert(ctx context.Context, u *UserModel) (status bool) {
	query, _, err := squirrel.
		Insert("users").
		Columns("username", "password").
		Values(u.username, u.password).
		ToSql()

	if err != nil {
		return false
	}

	_, err = m.Conn.Exec(query)
	if err != nil {
		return false
	}

	return true
}

func (m *UserRepository) Fetch(ctx context.Context, username string) (users []UserModel) {
	query, _, err := squirrel.
		Select("*").
		From("users").
		Where(squirrel.Eq{"username": username}).
		ToSql()

	if err != nil {
		return []UserModel{}
	}

	_, err = m.Conn.Query(query)
	if err != nil {
		return []UserModel{}
	}

	return []UserModel{}
}
