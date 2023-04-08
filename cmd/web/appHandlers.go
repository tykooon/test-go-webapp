package main

import (
	"net/http"

	"github.com/tykooon/test-go-webapp/pkg/messagedb"
)

func (app *app) homeHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		app.notFound(writer)
	} else {

		messageList := make([]*messagedb.Message, 0)
		messageList = append(messageList, app.messages.GetMessageById(1))
		messageList = append(messageList, app.messages.GetMessageById(2))
		app.templates.home.ExecuteTemplate(writer, "base.html", messageList)
	}
}

func (app *app) notFound(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}
