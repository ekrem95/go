package main

import (
	"database/sql"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func signup(c echo.Context) error {
	name, email, password := c.FormValue("name"), c.FormValue("email"), c.FormValue("password")

	var id int
	row := db.QueryRow(`select id from users where email=$1`, email)
	switch errDB := row.Scan(&id); errDB {
	case sql.ErrNoRows:
		hash, error := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if error != nil {
			panic(error)
		}

		if _, err = db.Exec(`insert into users (name, email, password) values ($1, $2, $3)`, name, email, hash); err != nil {
			panic(err)
		}
		setSession(c, name)
	case nil:
		return c.JSON(200, map[string]string{"msg": "User with same email address already exists."})
	default:
		c.JSON(200, map[string]string{"msg": "There was an error signing you up."})
		panic(err)
	}

	return c.JSON(200, map[string]bool{"done": true})
}

func login(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	var uname, dbPassword string

	err = db.QueryRow("select name, password from users where email=$1", email).Scan(&uname, &dbPassword)
	if err != nil {
		return c.JSON(200, map[string]string{"msg": "Internal server error."})
	}

	if err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password)); err != nil {
		return c.JSON(200, map[string]string{"msg": "Wrong email & password combination."})
	}

	setSession(c, uname)
	t := jwtAuth(uname)
	return c.JSON(200, map[string]string{"name": uname, "token": t})
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.JSON(200, map[string]string{"user": name})
}

func jwtAuth(username string) string {
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		panic(err)
	}
	return t
}
