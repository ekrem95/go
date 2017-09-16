package main

import "errors"

type user struct {
	Username string `json:"username"`
	Password string `json:"-"`
}

var userList = []user{
	user{Username: "user1", Password: "pass1"},
	user{Username: "user2", Password: "pass2"},
	user{Username: "user3", Password: "pass3"},
	user{Username: "user4", Password: "pass4"},
}

func registerNewUser(username, password string) (*user, error) {
	return nil, errors.New("placeholder error")
}

func isUsernameAvailable(username string) bool {
	return false
}
