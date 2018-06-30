package main

import (
	"encoding/json"
	"fmt"

	"github.com/couchbase/gocb"
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
	bucket, err := cluster.OpenBucket("eko", "123456")
	if err != nil {
		panic(err)
	}
	// fmt.Println(bucket)
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
}
