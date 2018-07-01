package main

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

func setSession(c echo.Context, value string) {
	s, _ := session.Get("session", c)
	s.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   600,
		HttpOnly: true,
	}
	s.Values["user"] = value
	s.Save(c.Request(), c.Response())
}

func getSession(c echo.Context) string {
	s, _ := session.Get("session", c)
	username := s.Values["user"]
	return username.(string)
}
