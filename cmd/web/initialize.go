package main

import (
	"github.com/tykooon/test-go-webapp/pkg/usersys"
)

// Sqlite constants
const DbFileName = "sql/messages.db"
const SqliteCreateQuery = `
DROP TABLE IF EXISTS Authors;
CREATE TABLE IF NOT EXISTS Authors( Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, Name TEXT NOT NULL, CreatedAt DATETIME);
DROP TABLE IF EXISTS Messages;
CREATE TABLE IF NOT EXISTS Messages( Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, AuthorId INTEGER NOT NULL, CreatedAt DATETIME, Content TEXT NOT NULL);`

// Mysql constants
//const DbSchemaName = "heroku_d869720877b7819"
//const MysqlConnectionString = "b20b350404308e:1cef0188@tcp(eu-cdbr-west-03.cleardb.net)/heroku_d869720877b7819"

const DbSchemaName = "test-go-webapp"
const MysqlConnectionString = "root:@tcp(127.0.0.1:3306)/test-go-webapp"

func init() {
	_ = usersys.User{ // TODO Delete this
		FirstName: "Alex",
		LastName:  "Tykoun",
	}
}
