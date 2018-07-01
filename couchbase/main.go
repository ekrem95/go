package main

import (
	"encoding/json"
	"fmt"

	"github.com/couchbase/gocb"
)

var (
	name = "default"
)

// Person ...
type Person struct {
	Firstname string `json:"firstname,omitempty`
	Lastname  string `json:"lastname,omitempty`
	Social    []Socialmedia
}

// Socialmedia ...
type Socialmedia struct {
	Title string `json:"title`
	Link  string `json:"link`
}

func main() {
	cluster, err := gocb.Connect("couchbase://127.0.0.1")
	if err != nil {
		panic(err)
	}
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "123456",
	})
	bucket, err := cluster.OpenBucket(name, "")
	if err != nil {
		panic(err)
	}
	bucket.Manager("", "").CreatePrimaryIndex("", true, false)

	var person Person

	// first param of upsert func is ID
	bucket.Upsert("ekrem", Person{
		Firstname: "Ekrem",
		Lastname:  "K",
		Social: []Socialmedia{
			{Title: "twitter", Link: "www"},
			{Title: "github", Link: "www"},
		},
	}, 0)

	bucket.Upsert("eko", Person{
		Firstname: "Ekrem",
		Lastname:  "K",
		Social: []Socialmedia{
			{Title: "twitter", Link: "www"},
			{Title: "github", Link: "www"},
		},
	}, 0)

	bucket.Get("ekrem", &person)
	jsonBytes, _ := json.Marshal(person)
	fmt.Println(string(jsonBytes))

	bucket.Get("eko", &person)
	jsonBytes, _ = json.Marshal(person)
	fmt.Println(string(jsonBytes))

	// Use query
	query := gocb.NewN1qlQuery("SELECT * FROM " + name + " WHERE $1 IN Social")
	rows, err := bucket.ExecuteN1qlQuery(query, []interface{}{Socialmedia{Title: "twitter", Link: "www"}})
	if err != nil {
		panic(err)
	}
	var row interface{}
	for rows.Next(&row) {
		fmt.Printf("Row: %v\n", row)
	}
}
