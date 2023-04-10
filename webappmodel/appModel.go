package webappmodel

import (
	"log"

	"github.com/tykooon/test-go-webapp/pkg/messagedb"
)

type App struct {
	Log       *log.Logger
	messages  *messagedb.MessageDB
	templates *templates
}

func NewApp(log *log.Logger, messages *messagedb.MessageDB) *App {
	return &App{
		Log:       log,
		messages:  messages,
		templates: initTemplates(),
	}
}
