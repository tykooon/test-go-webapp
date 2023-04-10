package dbprovider

import (
	"database/sql"
)

const (
	ERR_NO_CONNECTION = "No connection to database"
)

type DbProvider interface {
	OpenDB() (*sql.DB, error)
}
