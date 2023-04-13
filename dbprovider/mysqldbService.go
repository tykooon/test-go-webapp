package dbprovider

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tykooon/test-go-webapp/pkg/messagedb"
)

type mysqlDbService struct {
	db            *sql.DB
	connectionUrl string
	dbSchemaName  string
}

func NewMysqlDbService(connectionUrl string, schema string) *mysqlDbService {
	return &mysqlDbService{
		connectionUrl: connectionUrl,
		dbSchemaName:  schema,
	}
}

func (m *mysqlDbService) OpenDB() (err error) {
	m.db, err = sql.Open("mysql", m.connectionUrl+"?parseTime=true")
	if err != nil {
		return
	}
	rows, err := m.db.Query("SHOW TABLES IN `" + m.dbSchemaName + "`;")
	if err != nil {
		m.db.Close()
		return
	}
	defer rows.Close()
	if !rows.Next() {
		return fmt.Errorf("no tables found")
	}
	return err
}

func (m *mysqlDbService) DB() *sql.DB {
	return m.db
}

func (m *mysqlDbService) GetAuthorByName(name string) *messagedb.Author {
	author := messagedb.NewAuthor(name)
	row := m.db.QueryRow(`
		SELECT Id 
		FROM authors 
		WHERE name = ?;`,
		name)
	if row.Scan(&author.Id) != nil {
		return nil
	}
	return author
}

func (m *mysqlDbService) GetAuthorById(id int64) *messagedb.Author {
	author := &messagedb.Author{Id: id}
	row := m.db.QueryRow(`
		SELECT Name, CreatedAt
		FROM authors
		WHERE id =?;`,
		id)
	if row.Scan(&author.Name, &author.CreatedAt) != nil {
		return nil
	}
	return author
}

func (m *mysqlDbService) GetMessageById(id int64) *messagedb.Message {
	message := &messagedb.Message{Id: id}
	var authId int64
	row := m.db.QueryRow(`
		SELECT AuthorId, CreatedAt, Content
		FROM messages
		WHERE id =?;`,
		id)
	if err := row.Scan(&authId, &message.CreatedAt, &message.Content); err != nil {
		return nil
	}
	message.Author = *m.GetAuthorById(authId)
	return message
}

func (m *mysqlDbService) GetAllMessages() []*messagedb.Message {
	rows, err := m.db.Query(`
		SELECT AuthorId, CreatedAt, Content
		FROM messages
		ORDER BY CreatedAt DESC;`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	messages := make([]*messagedb.Message, 0)
	for rows.Next() {
		message := &messagedb.Message{}
		if err := rows.Scan(&message.Author.Id, &message.CreatedAt, &message.Content); err != nil {
			return nil
		}
		message.Author = *m.GetAuthorById(message.Author.Id)
		messages = append(messages, message)
	}
	return messages
}

func (m *mysqlDbService) InsertMessageWithAuthor(authorName string, content string) int64 {
	if authorName == "" || content == "" {
		return ErrWrongData
	}
	author := m.GetAuthorByName(authorName)
	if author == nil {
		author = messagedb.NewAuthor(authorName)
		author.Id = m.InsertAuthor(author)
		if author.Id == ErrAuthorNotCreated {
			return ErrAuthorNotCreated
		}
	}
	message := messagedb.NewMessage(author, content)
	messId := m.InsertMessage(message)
	if messId == ErrMessageNotCreated {
		return ErrMessageNotCreated
	}
	return messId
}

func (m *mysqlDbService) InsertAuthor(author *messagedb.Author) (id int64) {
	if author.Name == "" {
		return ErrWrongData
	}
	if m.GetAuthorByName(author.Name) != nil {
		return ErrEntryAlreadyExists
	}
	row, err := m.db.Exec(`
		INSERT INTO authors (Name, CreatedAt)
		VALUES (?, UTC_TIMESTAMP());`,
		author.Name)
	if err == nil {
		id, err = row.LastInsertId()
		if err == nil {
			author.Id = id
			return
		}
	}
	fmt.Println(err.Error())
	return ErrAuthorNotCreated
}

func (m *mysqlDbService) InsertMessage(message *messagedb.Message) (id int64) {
	if message == nil || message.Content == "" || message.Author.Id <= 0 {
		return ErrWrongData
	}
	row, err := m.db.Exec(`
		INSERT INTO messages (AuthorId, CreatedAt, Content)
		VALUES ( ? ,UTC_TIMESTAMP(), ? );`,
		message.Author.Id,
		message.Content)
	if err == nil {
		id, err = row.LastInsertId()
		if err == nil {
			message.Id = id
			return
		}
	}
	fmt.Println(err.Error())
	return ErrMessageNotCreated
}
