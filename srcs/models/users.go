package models

import (
	"database/sql"
	"errors"
)

//Use to store data that we retrieve from db
type User struct {
	Id int
	Username string
	Email string
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) GetByUsername(username string) (*User, error) {
	var user User
	//QueryRow => query just one row from the table
	row := m.DB.QueryRow("SELECT id, username, email FROM users WHERE username = $1", username)
	err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Email,
	)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, err
}
