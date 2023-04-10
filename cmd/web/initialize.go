package main

import (
	"github.com/tykooon/test-go-webapp/pkg/usersys"
)

func init() {
	_ = usersys.User{ // TODO Delete this
		FirstName: "Alex",
		LastName:  "Tykoun",
	}
}
