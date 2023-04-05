package main

import (
	"log"
	"net/http"
	"os"

	"github.com/tykooon/test-go-webapp/pkg/usersys"
)

func main() {

	logger := log.New(os.Stdout, "APP LOG :", log.Lshortfile)

	port := os.Getenv("PORT")
	if port == "" {
		logger.Fatal("Error: port must be set")
	}

	u1 := usersys.User{
		FirstName: "Alex",
		LastName:  "Tykoun",
	}

	app := app{
		user:      u1,
		log:       logger,
		templates: initTemplates(),
	}

	err := http.ListenAndServe(":"+port, app.routes())
	app.log.Fatal(err)
}
