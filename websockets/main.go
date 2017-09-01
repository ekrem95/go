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
		}
		fmt.Println("Clien subscribed")

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				panic(err)
			}
			if string(msg) == "ping" {
				fmt.Println("ping")
				time.Sleep(2 * time.Second)
				err = conn.WriteMessage(msgType, []byte("pong"))
				if err != nil {
					panic(err)
				}
			} else {
				conn.Close()
				fmt.Println(string(msg))
				return
			}
		}
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(index))
	})
	http.ListenAndServe(":3000", nil)
}
