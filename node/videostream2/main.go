package main

import (
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/googollee/go-socket.io"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func main() {
	// Echo instance
	e := echo.New()

	res := client.Del("test_online_users")
	if res.Err() != nil {
		panic(res.Err())
	}

	server, err := websocket()
	if err != nil {
		log.Fatal(err)
	}

	// Middleware
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/streams", smembers)
	e.GET("/socket.io/", echo.WrapHandler(server))
	e.POST("/socket.io/", echo.WrapHandler(server))

	e.File("/bundle.js", "public/bundle.js")
	e.File("*", "template/index.html")

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// Handler
func smembers(c echo.Context) error {
	members, _ := client.SMembers("test_online_users").Result()
	return c.JSON(http.StatusOK, members)
}

func websocket() (*socketio.Server, error) {
	server, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}
	var user string

	server.On("connection", func(so socketio.Socket) {
		log.Println("new client connected")

		so.Join("new_stream")
		so.Join("stream")

		so.On("new_stream", func(uname string) {
			user = uname
			log.Printf("%s connected", uname)
			client.SAdd("test_online_users", uname)
		})
		so.On("stream", func(data []string) {
			so.BroadcastTo("stream", "dist"+data[1], data[0])
		})
		so.On("disconnection", func() {
			log.Printf("%s disconnected", user)
			client.SRem("test_online_users", user)
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	return server, nil
}
