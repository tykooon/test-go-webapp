package main

import (
	"log"
	"net/http"
	"os"

	"github.com/tykooon/test-go-webapp/dbprovider"
	"github.com/tykooon/test-go-webapp/pkg/messagedb"
	"github.com/tykooon/test-go-webapp/webappmodel"
)

const DbFileName = "sql/messages.db"

func main() {
	logger := log.New(os.Stdout, "APP LOG :", log.Lshortfile)

	port := GetPortFromEnviron()

	dbservice := dbprovider.NewSqliteDbService(DbFileName, createQuery)

	db, err := dbservice.OpenDB()
	if err != nil {
		logger.Fatal("Sorry. Failed to open Database... ", err.Error())
		return
	}
	defer db.Close()

	mdb := messagedb.NewMessageDB(db)
	app := webappmodel.NewApp(logger, mdb)

	err = http.ListenAndServe(":"+port, app.Routes())
	app.Log.Fatal(err)
}

func GetPortFromEnviron() (port string) {
	port = os.Getenv("PORT")
	if port == "" {
		port = "5000"
		// TODO Only for local test
		//             Replace with
		//logger.Fatal("Error: port must be set")
	}
	return
}
