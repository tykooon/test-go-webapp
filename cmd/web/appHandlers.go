package main

import (
	"net/http"
)

func (app *app) homeHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		app.notFound(writer)
	} else {
		app.templates.home.ExecuteTemplate(writer, "base.html", app.user)
	}
}

func (app *app) notFound(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}
