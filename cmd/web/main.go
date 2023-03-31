package main

import (
	"fmt"

	"github.com/tykooon/test-go-webapp/pkg/usersys"
)

func main() {

	u1 := usersys.User{
		FirstName: "Alex",
		LastName:  "Tykoun",
	}

	fmt.Println(u1)

}
