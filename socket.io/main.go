package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/googollee/go-socket.io"
)

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")
		so.Join("chat")
		so.On("chat message", func(msg string) {
			log.Println("emit:", so.Emit("chat message", msg))
			so.BroadcastTo("chat", "chat message", msg)
		})
		so.On("disconnection", func() {
			log.Println("on disconnect")
		})

		so.On("time", func(time string) string {
			log.Println(time)

			// return time[0:10]
			return strings.Replace(time[0:10], "-", " ", -1)
		})
	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	fs := http.FileServer(http.Dir("./"))

	http.Handle("/socket.io/", server)
	http.Handle("/", fs)
	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
