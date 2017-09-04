package main

import (
	"database/sql"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func signup(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")

	var id int
	row := db.QueryRow(`select id from users where email=$1`, email)
	switch err := row.Scan(&id); err {
	case sql.ErrNoRows:
		hash, error := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if error != nil {
			panic(error)
		}

		_, err = db.Exec(`insert into users (name, email, password) values ($1, $2, $3)`, name, email, hash)
		if err != nil {
			panic(err)
		}
	case nil:
		return c.JSON(200, map[string]string{
			"msg": "User with same email address already exists.",
		})
	default:
		c.JSON(200, map[string]string{
			"msg": "There was an error signing you up.",
		})
		panic(err)
	}

	return c.JSON(200, map[string]bool{
		"done": true,
	})
}
