package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	database = "default"
	host     = "localhost"
	password = "password"
	port     = 5432
	user     = "ekrem"
)

var db *sql.DB
var err error

func psql() {
	info := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, database)

	db, err = sql.Open("postgres", info)
	if err != nil {
		panic(err)
	}
	// defer db.Close()

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")

	stmt :=
		`
		CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR (50) NOT NULL,
		email VARCHAR (150) UNIQUE NOT NULL,
		password VARCHAR (150) NOT NULL
		);
		`

	if _, err = db.Exec(stmt); err != nil {
		panic(err)
	}
}

// User type
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
