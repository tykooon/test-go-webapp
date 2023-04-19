package dbprovider

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/tykooon/test-go-webapp/pkg/messagedb"
)

type pgsqlDbService struct {
	db            *sql.DB
	connectionUrl string
	dbSchemaName  string
}

func NewPgsqlDbService(connectionUrl string, schema string) *pgsqlDbService {
	return &pgsqlDbService{
		connectionUrl: connectionUrl,
		dbSchemaName:  schema,
	}
}

func (p *pgsqlDbService) OpenDB() (err error) {
	p.db, err = sql.Open("postgres", p.connectionUrl) // +"?parseTime=true")
	if err != nil {
		return
	}
	err = p.db.Ping()
	if err != nil {
		p.db.Close()
	}
	return
}

func (p *pgsqlDbService) DB() *sql.DB {
	return p.db
}

func (p *pgsqlDbService) GetAuthorByName(name string) *messagedb.Author {
	author := messagedb.NewAuthor(name)
	row := p.db.QueryRow(`
		SELECT Id 
		FROM dbo.authors 
		WHERE name = $1;`,
		name)
	if row.Scan(&author.Id) != nil {
		return nil
	}
	return author
}

func (p *pgsqlDbService) GetAuthorById(id int64) *messagedb.Author {
	author := &messagedb.Author{Id: id}
	row := p.db.QueryRow(`
		SELECT Name, CreatedAt
		FROM dbo.authors
		WHERE id =$1;`,
		id)
	if row.Scan(&author.Name, &author.CreatedAt) != nil {
		return nil
	}
	return author
}

func (p *pgsqlDbService) GetMessageById(id int64) *messagedb.Message {
	message := &messagedb.Message{Id: id}
	var authId int64
	row := p.db.QueryRow(`
		SELECT AuthorId, CreatedAt, Content
		FROM dbo.messages
		WHERE id =$1;`,
		id)
	if err := row.Scan(&authId, &message.CreatedAt, &message.Content); err != nil {
		return nil
	}
	message.Author = *p.GetAuthorById(authId)
	return message
}

func (p *pgsqlDbService) GetAllMessages() []*messagedb.Message {
	rows, err := p.db.Query(`
		SELECT AuthorId, CreatedAt, Content
		FROM dbo.messages
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
		message.Author = *p.GetAuthorById(message.Author.Id)
		messages = append(messages, message)
	}
	return messages
}

func (p *pgsqlDbService) InsertMessageWithAuthor(authorName string, content string) int64 {
	if authorName == "" || content == "" {
		return ErrWrongData
	}
	author := p.GetAuthorByName(authorName)
	if author == nil {
		author = messagedb.NewAuthor(authorName)
		author.Id = p.InsertAuthor(author)
		if author.Id == ErrAuthorNotCreated {
			return ErrAuthorNotCreated
		}
	}
	message := messagedb.NewMessage(author, content)
	messId := p.InsertMessage(message)
	if messId == ErrMessageNotCreated {
		return ErrMessageNotCreated
	}
	return messId
}

func (p *pgsqlDbService) InsertAuthor(author *messagedb.Author) (id int64) {
	if author.Name == "" {
		return ErrWrongData
	}
	if p.GetAuthorByName(author.Name) != nil {
		return ErrEntryAlreadyExists
	}
	row := p.db.QueryRow(`
		INSERT INTO dbo.authors (Name, CreatedAt)
		VALUES ($1, now())
		RETURNING id`,
		author.Name)
	err := row.Scan(&id)
	if err == nil {
		author.Id = id
		return
	}
	fmt.Println(err.Error())
	return ErrAuthorNotCreated
}

func (p *pgsqlDbService) InsertMessage(message *messagedb.Message) (id int64) {
	if message == nil || message.Content == "" || message.Author.Id <= 0 {
		return ErrWrongData
	}
	row := p.db.QueryRow(`
		INSERT INTO dbo.messages (AuthorId, CreatedAt, Content)
		VALUES ( $1, now(), $2 )
		RETURNING id`,
		message.Author.Id,
		message.Content)
	err := row.Scan(&id)
	if err == nil {
		message.Id = id
		return
	}
	fmt.Println(err.Error())
	return ErrMessageNotCreated
}
