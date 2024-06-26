package models

import (
	"database/sql"
	"errors"
)

//Use to store data that we retrieve from db
type User struct {
	Id int
	// Username string
	Username sql.NullString
	Email string
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) GetListWithUsername() ([]User, error) {
	rows, err := m.DB.Query("SELECT id, username, email FROM users WHERE username IS NOT NULL")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows){
		_ = rows.Close()
	}(rows)

	var users []User

	//Loop through each rows
	for rows.Next() {
		var user User
		err = rows.Scan(
			&user.Id,
			&user.Username,
			&user.Email,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
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

func (m *UserModel) GetByEmail(email string) (*User, error) {
	var user User
	//QueryRow => query just one row from the table
	row := m.DB.QueryRow("SELECT id, username, email FROM users WHERE email = $1", email)
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
