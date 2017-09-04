package main

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

func setSessionUser(c echo.Context, value string) {
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   600,
		HttpOnly: true,
	}
	sess.Values["user"] = value
	sess.Save(c.Request(), c.Response())
}

func getSessionUser(c echo.Context) string {
	sess, _ := session.Get("session", c)
	username := sess.Values["user"]
	return username.(string)
}
