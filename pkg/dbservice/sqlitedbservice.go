package dbservice

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

type sqliteDbService struct {
	dbFileName  string
	createQuery string
}

func NewSqliteDbService(dbFileName string, createQuery string) *sqliteDbService {
	return &sqliteDbService{
		dbFileName:  dbFileName,
		createQuery: createQuery,
	}
}

func (s *sqliteDbService) OpenDB() (db *sql.DB, err error) {
	_, err = os.Stat(s.dbFileName)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		file, err := os.Create(s.dbFileName)
		if err != nil {
			return nil, err
		}
		file.Close()
	}
	db, err = sql.Open("sqlite", s.dbFileName)
	if err == nil {
		err = db.Ping()
		if err == nil {
			_, err = db.Exec(s.createQuery)
			return db, err
		}
	}
	return nil, err
}
