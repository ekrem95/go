package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

// Person type
type Person struct {
	Name string
	Age  int
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	indexFile, err := os.Open("./html/index.html")
	if err != nil {
		panic(err)
	}
	index, err := ioutil.ReadAll(indexFile)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err)
			return
		}
		fmt.Println("Clien subscribed")

		myPerson := Person{Name: "Bill", Age: 0}
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(index))
	})
	http.ListenAndServe(":4000", nil)
}
