package main

import (
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"
	"github.com/tykooon/test-go-webapp/dbprovider"
	"github.com/tykooon/test-go-webapp/webappmodel"
)

func main() {
	logger := log.New(os.Stdout, "APP LOG :", log.Lshortfile)

	config := viper.New()
	config.SetConfigName("appsettings")
	config.SetConfigType("yaml")
	config.AddConfigPath(".")
	err := config.ReadInConfig()
	if err != nil {
		logger.Fatal("Missing configuration.")
	}

	dbType := config.GetString("dbType")
	dbParams := config.GetStringMapString("dbProfiles." + dbType)

	dbservice, err := dbprovider.NewDbService(dbType, dbParams)
	if err != nil {
		logger.Fatal(err)
	}

	err = dbservice.OpenDB()
	if err != nil {
		logger.Fatal("Sorry. Failed to open Database... ", err.Error())
		return
	}
	defer dbprovider.CloseDB(dbservice)

	app := webappmodel.NewApp(logger, dbservice)

	port := GetPortFromEnviron()

	err = http.ListenAndServe(":"+port, app.Routes())
	app.Log.Fatal(err)
}

func GetPortFromEnviron() (port string) {
	port = os.Getenv("PORT")
	if port == "" {
		port = "5000"
		// TODO `port = "5000"` --- only for local test
		//            In realize replace with
		//	logger.Fatal("Error: port must be set")
	}
	return
}
