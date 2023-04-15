package main

import (
	"github.com/tykooon/test-go-webapp/pkg/usersys"
)

// Mysql constants
//const DbSchemaName = "heroku_d869720877b7819"
//const MysqlConnectionString = "b20b350404308e:1cef0188@tcp(eu-cdbr-west-03.cleardb.net)/heroku_d869720877b7819"

func init() {
	_ = usersys.User{ // TODO Delete this
		FirstName: "Alex",
		LastName:  "Tykoun",
	}
}
