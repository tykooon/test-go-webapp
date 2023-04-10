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
	if err = s.ExistsDB(); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		err = s.CreateDB()
		if err != nil {
			return nil, err
		}
	}
	db, err = sql.Open("sqlite", s.dbFileName)
	if err == nil {
		err = db.Ping()
		if err == nil {
			return db, err
		}
	}
	return nil, err
}

func (s *sqliteDbService) ExistsDB() error {
	_, err := os.Stat(s.dbFileName)
	if err != nil {
		return err
	}
	return nil
}

func (s *sqliteDbService) CreateDB() error {
	file, err := os.Create(s.dbFileName)
	if err != nil {
		return err
	}
	file.Close()
	db, err := sql.Open("sqlite", s.dbFileName)
	defer db.Close()
	if err == nil {
		err = db.Ping()
		if err == nil {
			_, err = db.Exec(s.createQuery)
			return err
		}
	}
	return err
}
