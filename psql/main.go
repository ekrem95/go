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

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL,
		age INT,
		first_name VARCHAR(255),
		last_name VARCHAR(255),
		email TEXT
	  );`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO users (age, email, first_name, last_name)
	VALUES (30, 'eko@eko.io', 'Ekrem', 'Karatas'); `)
	if err != nil {
		panic(err)
	}

	rows, err := db.Query(`select * from users`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var id int
		var age int
		var first_name string
		var last_name string
		var email string
		err = rows.Scan(&id, &age, &first_name, &last_name, &email)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("%3v | %8v | %6v | %16v | %6v\n", id, age, first_name, last_name, email)
	}
}
