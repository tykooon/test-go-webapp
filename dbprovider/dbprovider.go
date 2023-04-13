package dbprovider

import (
	"database/sql"

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
