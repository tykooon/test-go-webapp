package main

import (
	"log"
	"net/http"

	"github.com/tykooon/test-go-webapp/pkg/messagedb"
	"github.com/tykooon/test-go-webapp/pkg/usersys"
)

type app struct {
	user      usersys.User
	log       *log.Logger
	messages  *messagedb.MessageDB
	templates *templates
}

func (app *app) routes() *http.ServeMux {

	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(app.homeHandler))
	return mux
}
