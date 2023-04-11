package dbprovider

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type mysqlDbService struct {
	connectionUrl string
	createQuery   string
}

func NewMysqlDbService(connectionUrl string, createQuery string) *mysqlDbService {
	return &mysqlDbService{
		connectionUrl: connectionUrl,
		createQuery:   createQuery,
	}
}

func (m *mysqlDbService) OpenDB() (db *sql.DB, err error) {
	// TODO: recode
	if _, err = os.Stat(m.connectionUrl); err != nil {
		if os.IsNotExist(err) {
			db, err = m.OpenAndPing()
			if err == nil {
				_, err = db.Exec(m.createQuery)
			}
		}
	} else {
		db, err = m.OpenAndPing()
	}
	return db, err
}

func (m *mysqlDbService) OpenAndPing() (db *sql.DB, err error) {
	//TODO recode
	db, err = sql.Open("mysql", m.connectionUrl)
	if err == nil {
		err = db.Ping()
	}
	return db, err
}
