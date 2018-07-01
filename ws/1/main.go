package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	file, err := os.Open("./html/index.html")
	if err != nil {
		panic(err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err)
		}
		fmt.Println("Client subscribed")

		for {
			ch, msg, err := conn.ReadMessage()
			if err != nil {
				panic(err)
			}
			fmt.Println(string(msg))

			if string(msg) == "ping" {
				time.Sleep(2 * time.Second)
				if err = conn.WriteMessage(ch, []byte("pong")); err != nil {
					panic(err)
				}
			} else {
				conn.Close()
				fmt.Println("Connection closed")
				break
			}
		}
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(content))
	})
	http.ListenAndServe(":8080", nil)
}
