package main

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
)

func main() {
	psql()
	e := echo.New()
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowCredentials: true,
	}))
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	e.POST("/signup", signup)
	e.POST("/login", login)

	e.GET("/", accessible)

	r := e.Group("/restricted")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("", restricted)

	e.Logger.Fatal(e.Start(":1323"))
}
