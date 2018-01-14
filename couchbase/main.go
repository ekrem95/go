package main

import (
	"encoding/json"
	"fmt"

	"github.com/couchbase/gocb"
)

type Person struct {
	ID        string `json:"id,omitempty`
	Firstname string `json:"firstname,omitempty`
	Lastname  string `json:"lastname,omitempty`
	Social    []Socialmedia
}

type Socialmedia struct {
	Title string `json:"title`
	Link  string `json:"link`
}

func main() {
	cluster, err := gocb.Connect("couchbase://127.0.0.1")
	if err != nil {
		fmt.Println(err)
	}
	bucket, err := cluster.OpenBucket("eko", "123456")
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(bucket)
	var person Person

	// first param of upsert func is ID
	bucket.Upsert("ekrem", Person{
		ID:        "1",
		Firstname: "Ekrem",
		Lastname:  "K",
		Social: []Socialmedia{
			{Title: "twitter", Link: "www"},
			{Title: "github", Link: "www"},
		},
	}, 0)

	bucket.Upsert("eko", Person{
		ID:        "2",
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
