package models

import (
	"errors"

	"example.com/rest-api/db"
	"example.com/rest-api/utils"
)

type User struct {
	ID    int64
	Email string `binding:"required"`
	Pswd  string `binding:"required"`
}

func (u User) Save() error {
	query := "INSERT INTO users(email, pswd) VALUES (?, ?)"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()
	// encoding raw password into hashed one
	hashedPswd, err := utils.HashPassword(u.Pswd)

	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Email, hashedPswd)

	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()
	u.ID = userId
	return err
}

func GetUsers() ([]User, error) {
	query := "SELECT * FROM users"
	rows, err := db.DB.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.Pswd)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	return users, nil
}

func (u *User) UserValidation() error {
	query := "SELECT id, pswd FROM users WHERE email = ?"
	row := db.DB.QueryRow(query, u.Email)

	var dbStoredPswd string
	err := row.Scan(&u.ID, &dbStoredPswd)

	if err != nil {
		return errors.New("email or password is incorrect") // here checks email
	}

	pswdIsValid := utils.CompareHashAndPassword(u.Pswd, dbStoredPswd)

	if !pswdIsValid {
		return errors.New("email or password is incorrect") // here checks password
	}
	return nil
}
