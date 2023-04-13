package main

import (
	"log"
	"net/http"
	"os"

	"github.com/tykooon/test-go-webapp/dbprovider"
	"github.com/tykooon/test-go-webapp/webappmodel"
)

func main() {
	logger := log.New(os.Stdout, "APP LOG :", log.Lshortfile)

	port := GetPortFromEnviron()

	dbservice := dbprovider.NewMysqlDbService(MysqlConnectionString, DbSchemaName)
	//dbservice := dbprovider.NewSqliteDbService(DbFileName, SqliteCreateQuery)

	err := dbservice.OpenDB()
	if err != nil {
		logger.Fatal("Sorry. Failed to open Database... ", err.Error())
		return
	}
	defer dbprovider.CloseDB(dbservice)

	app := webappmodel.NewApp(logger, dbservice)

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
