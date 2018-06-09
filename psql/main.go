package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	DB_USER     = "username"
	DB_PASSWORD = "password"
	DB_NAME     = "default"
)

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL,
		age INT,
		first_name VARCHAR(255),
		last_name VARCHAR(255),
		email TEXT
	  );`); err != nil {
		panic(err)
	}
	if _, err = db.Exec(`INSERT INTO users (age, email, first_name, last_name)
	VALUES (30, 'eko@eko.io', 'Ekrem', 'Karatas'); `); err != nil {
		panic(err)
	}

	rows, err := db.Query(`select * from users`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var id, age int
		var fname, lname, email string
		if err = rows.Scan(&id, &age, &fname, &lname, &email); err != nil {
			fmt.Println(err)
		}

		fmt.Printf("%3v | %8v | %6v | %16v | %6v\n", id, age, fname, lname, email)
	}
}
