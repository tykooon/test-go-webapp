package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/tykooon/test-go-webapp/pkg/dbservice"
	"github.com/tykooon/test-go-webapp/pkg/messagedb"
	"github.com/tykooon/test-go-webapp/pkg/usersys"
)

const DbFileName = "sql/messages.db"

func main() {

	logger := log.New(os.Stdout, "APP LOG :", log.Lshortfile)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000" // Only for local test
		//logger.Fatal("Error: port must be set")
	}

	u1 := usersys.User{
		FirstName: "Alex",
		LastName:  "Tykoun",
	}

	//
	dbservice := dbservice.NewSqliteDbService(DbFileName, createQuery)

	db, err := dbservice.OpenDB()
	if err != nil {
		logger.Fatal("Sorry. Failed to open Database... ", err.Error())
		return
	}
	defer db.Close()

	app := app{
		user:      u1,
		messages:  messagedb.NewMessageDB(db),
		log:       logger,
		templates: initTemplates(),
	}

	auth1 := &messagedb.Author{Name: "Alex T."}
	fmt.Println(app.messages.InsertAuthor(auth1))
	fmt.Println(auth1)
	mess1 := messagedb.NewMessage(auth1, "Some message")
	fmt.Println(app.messages.InsertMessage(mess1))

	fmt.Println(app.messages.Insert("Uliana", "Hi there!"))

	fmt.Println(app.messages.Insert("Uliana", "Second message"))

	fmt.Println(app.messages.GetAuthorById(1))
	fmt.Println(app.messages.GetAuthorByName("Uliana"))
	fmt.Println(app.messages.GetMessageById(2))

	err = http.ListenAndServe(":"+port, app.routes())
	app.log.Fatal(err)
}

const createQuery string = `
DROP TABLE IF EXISTS Authors;
CREATE TABLE IF NOT EXISTS Authors(
	Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	Name TEXT NOT NULL,
	CreatedAt TIMESPAMP);
DROP TABLE IF EXISTS Messages;
CREATE TABLE IF NOT EXISTS Messages(
	Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	AuthorId INTEGER NOT NULL,			
	CreatedAt TIMESPAMP,
	Content TEXT NOT NULL);`
