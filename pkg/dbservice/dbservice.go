package dbservice

import (
	"database/sql"
)

const (
	ErrNoConnection = "No connection with database"
)

type DbProvider interface {
	OpenDB() (*sql.DB, error)
}
