package dbprovider

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
	if _, err = os.Stat(s.dbFileName); err != nil {
		if os.IsNotExist(err) {
			err = s.CreateDBFile()
			if err == nil {
				db, err = s.OpenAndPing()
				if err == nil {
					_, err = db.Exec(s.createQuery)
				}
			}
		}
	} else {
		db, err = s.OpenAndPing()
	}
	return db, err
}

func (s *sqliteDbService) CreateDBFile() error {
	file, err := os.Create(s.dbFileName)
	if err == nil {
		file.Close()
	}
	return err
}

func (s *sqliteDbService) OpenAndPing() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite", s.dbFileName)
	if err == nil {
		err = db.Ping()
	}
	return db, err
}
