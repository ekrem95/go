package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
)

func main() {
	r := gin.Default()
	r.StaticFS("/io", http.Dir("./node_modules/socket.io-client/dist/"))
	r.LoadHTMLGlob("public/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.GET("/video", func(c *gin.Context) {
		c.HTML(http.StatusOK, "view.html", gin.H{})
	})

	// socketio
	server, err := websocket()
	if err != nil {
		log.Fatal(err)
	}

	// socketio
	r.GET("/socket.io/", gin.WrapH(server))
	r.POST("/socket.io/", gin.WrapH(server))

	r.Run(":3000") // listen and serve on 0.0.0.0:8080
}

func websocket() (*socketio.Server, error) {
	server, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}
	server.On("connection", func(so socketio.Socket) {
		log.Println("new client connected")

		so.Join("stream")

		so.On("stream", func(data string) {
			so.BroadcastTo("stream", "stream", data)
		})
		so.On("disconnection", func() {
			log.Println("client disconnected")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	return server, nil
}
