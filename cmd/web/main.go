package main

import (
	"log"
	"net/http"
	"os"

	"github.com/tykooon/test-go-webapp/pkg/usersys"
)

func main() {

	u1 := usersys.User{
		FirstName: "Alex",
		LastName:  "Tykoun",
	}

	logger := log.New(os.Stdout, "APP LOG :", log.Lshortfile)

	app := app{
		user:      u1,
		log:       logger,
		templates: initTemplates(),
	}

	err := http.ListenAndServe(":5000", app.routes())
	app.log.Fatal(err)
}
