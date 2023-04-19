package dbprovider

import (
	"database/sql"
	"fmt"

	"github.com/tykooon/test-go-webapp/pkg/messagedb"
)

const (
	ERR_NO_CONNECTION = "No connection to database"
)

const (
	ErrWrongData int64 = -(iota + 1)
	ErrEntryAlreadyExists
	ErrEntryNotFound
	ErrAuthorNotCreated
	ErrMessageNotCreated
)

type DbProvider interface {
	OpenDB() error
	DB() *sql.DB

	GetAuthorByName(name string) *messagedb.Author
	GetAuthorById(id int64) *messagedb.Author
	GetMessageById(id int64) *messagedb.Message
	GetAllMessages() []*messagedb.Message
	InsertMessageWithAuthor(authorName string, content string) int64
	InsertAuthor(author *messagedb.Author) (id int64)
	InsertMessage(message *messagedb.Message) (id int64)
}

func CloseDB(p DbProvider) error {
	if p == nil || p.DB() == nil {
		return nil
	}
	return p.DB().Close()
}

func NewDbService(serviceName string, paramDict map[string]string) (DbProvider, error) {
	switch serviceName {
	case "sqlite":
		arg1, ok1 := paramDict["dbfilename"]
		arg2, ok2 := paramDict["createquery"]
		fmt.Println(paramDict, arg1, arg2)
		if ok1 && ok2 {
			return NewSqliteDbService(arg1, arg2), nil
		}
	case "mysql":
		arg1, ok1 := paramDict["connectionurl"]
		arg2, ok2 := paramDict["dbschema"]
		if ok1 && ok2 {
			return NewMysqlDbService(arg1, arg2), nil
		}
	case "postgres":
		arg1, ok1 := paramDict["connectionurl"]
		arg2, ok2 := paramDict["dbschema"]
		if ok1 && ok2 {
			return NewPgsqlDbService(arg1, arg2), nil
		}
	default:
		return nil, fmt.Errorf("bad service name: %s", serviceName)
	}
	return nil, fmt.Errorf("bad config data for: %s", serviceName)
}
