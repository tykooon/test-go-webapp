package webappmodel

import (
	"fmt"
	"net/http"

	"github.com/tykooon/test-go-webapp/pkg/messagedb"
)

func (app *App) HomeHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		app.notFound(writer)
	} else {
		app.templates.home.ExecuteTemplate(writer, "base.html", nil)
	}
}

func (app *App) MessagesHandler(writer http.ResponseWriter, request *http.Request) {
	messageList := app.messages.GetAllMessages()
	app.templates.messages.ExecuteTemplate(writer, "base.html", messageList)
}

func (app *App) SendHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		app.templates.send.ExecuteTemplate(writer, "base.html", nil)
		return
	} else if request.Method == http.MethodPost {
		message := &messagedb.Message{}
		message.Content = request.PostFormValue("message")
		authName := request.PostFormValue("author")
		auth := app.messages.GetAuthorByName(authName)
		if auth == nil {
			auth = messagedb.NewAuthor(authName)
			authId := app.messages.InsertAuthor(auth)
			if authId <= 0 {
				fmt.Println(authId)
				app.serverError(writer)
				return
			}
		}
		message.Author = *auth
		messageId := app.messages.InsertMessage(message)
		if messageId <= 0 {
			app.serverError(writer)
			return
		}
		http.Redirect(writer, request, "/messages", http.StatusTemporaryRedirect)
		return
	}
	app.methodNotAllowed(writer)
}

func (app *App) notFound(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (app *App) serverError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *App) methodNotAllowed(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
