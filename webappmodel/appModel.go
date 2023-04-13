package webappmodel

import (
	"log"

	"github.com/tykooon/test-go-webapp/dbprovider"
)

type App struct {
	Log       *log.Logger
	messages  dbprovider.DbProvider
	templates *templates
}

func NewApp(log *log.Logger, messages dbprovider.DbProvider) *App {
	return &App{
		Log:       log,
		messages:  messages,
		templates: initTemplates(),
	}
}
