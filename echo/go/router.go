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
	switch errDB := row.Scan(&id); errDB {
	case sql.ErrNoRows:
		hash, error := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if error != nil {
			panic(error)
		}

		_, err = db.Exec(`insert into users (name, email, password) values ($1, $2, $3)`, name, email, hash)
		if err != nil {
			panic(err)
		}

		setSessionUser(c, name)
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

func login(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	var dbUsername string
	var dbPassword string

	err = db.QueryRow("select name, password from users where email=$1", email).Scan(&dbUsername, &dbPassword)
	if err != nil {
		return c.JSON(200, map[string]string{
			"msg": "Internal server error.",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	if err != nil {
		return c.JSON(200, map[string]string{
			"msg": "Wrong email & password combination.",
		})
	}

	setSessionUser(c, dbUsername)

	return c.JSON(200, map[string]string{
		"name": dbUsername,
	})
}
