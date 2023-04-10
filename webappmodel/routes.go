package webappmodel

import (
	"net/http"
)

func (app *App) Routes() *http.ServeMux {

	mux := http.NewServeMux()

	static := http.FileServer(http.Dir("./wwwroot"))

	mux.Handle("/files/", http.StripPrefix("/files", static))
	mux.Handle("/", http.HandlerFunc(app.HomeHandler))
	mux.Handle("/messages", http.HandlerFunc(app.MessagesHandler))
	mux.Handle("/send", http.HandlerFunc(app.SendHandler))
	return mux
}
