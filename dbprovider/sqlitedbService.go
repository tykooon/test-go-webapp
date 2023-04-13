package dbprovider

import (
	"database/sql"
	"os"

	"github.com/tykooon/test-go-webapp/pkg/messagedb"
	_ "modernc.org/sqlite"
)

type sqliteDbService struct {
	db          *sql.DB
	dbFileName  string
	createQuery string
}

func NewSqliteDbService(dbFileName string, createQuery string) *sqliteDbService {
	return &sqliteDbService{
		dbFileName:  dbFileName,
		createQuery: createQuery,
	}
}

func (s *sqliteDbService) OpenDB() (err error) {
	if _, err = os.Stat(s.dbFileName); err != nil {
		if os.IsNotExist(err) {
			err = s.CreateDBFile()
			if err == nil {
				err = s.OpenAndPing()
				if err == nil {
					_, err = s.db.Exec(s.createQuery)
				}
			}
		}
	} else {
		err = s.OpenAndPing()
	}
	return err
}

func (s *sqliteDbService) CreateDBFile() error {
	file, err := os.Create(s.dbFileName)
	if err == nil {
		file.Close()
	}
	return err
}

func (s *sqliteDbService) OpenAndPing() (err error) {
	s.db, err = sql.Open("sqlite", s.dbFileName+"?parseTime=true")
	if err == nil {
		err = s.db.Ping()
	}
	return err
}

func (s *sqliteDbService) DB() *sql.DB {
	return s.db
}

func (s *sqliteDbService) GetAuthorByName(name string) *messagedb.Author {
	author := messagedb.NewAuthor(name)
	row := s.db.QueryRow(`
		SELECT Id
		FROM authors
		WHERE name =?;`,
		name)
	if row.Scan(&author.Id) != nil {
		return nil
	}
	return author
}

func (s *sqliteDbService) GetAuthorById(id int64) *messagedb.Author {
	author := &messagedb.Author{Id: id}
	row := s.db.QueryRow(`
		SELECT Name, CreatedAt
		FROM authors
		WHERE id =?;`,
		id)
	if row.Scan(&author.Name, &author.CreatedAt) != nil {
		return nil
	}
	return author
}

func (s *sqliteDbService) GetMessageById(id int64) *messagedb.Message {
	message := &messagedb.Message{Id: id}
	var authId int64
	row := s.db.QueryRow(`
		SELECT AuthorId, CreatedAt, Content 
		FROM messages 
		WHERE id =?;`,
		id)
	if err := row.Scan(&authId, &message.CreatedAt, &message.Content); err != nil {
		return nil
	}
	message.Author = *s.GetAuthorById(authId)
	return message
}

func (s *sqliteDbService) GetAllMessages() []*messagedb.Message {
	rows, err := s.db.Query(`
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
		message.Author = *s.GetAuthorById(message.Author.Id)
		messages = append(messages, message)
	}
	return messages
}

func (s *sqliteDbService) InsertMessageWithAuthor(authorName string, content string) int64 {
	if authorName == "" || content == "" {
		return ErrWrongData
	}
	author := s.GetAuthorByName(authorName)
	if author == nil {
		author = messagedb.NewAuthor(authorName)
		author.Id = s.InsertAuthor(author)
		if author.Id == ErrAuthorNotCreated {
			return ErrAuthorNotCreated
		}
	}
	message := messagedb.NewMessage(author, content)
	messId := s.InsertMessage(message)
	if messId == ErrMessageNotCreated {
		return ErrMessageNotCreated
	}
	return messId
}

func (s *sqliteDbService) InsertAuthor(author *messagedb.Author) (id int64) {
	if author.Name == "" {
		return ErrWrongData
	}
	if s.GetAuthorByName(author.Name) != nil {
		return ErrEntryAlreadyExists
	}
	row := s.db.QueryRow(`
		INSERT INTO authors (Name, CreatedAt)
		VALUES (?, DATETIME('now','utc'))
		RETURNING Id`,
		author.Name)
	if row.Scan(&id) != nil {
		return ErrAuthorNotCreated
	}
	author.Id = id
	return
}

func (s *sqliteDbService) InsertMessage(message *messagedb.Message) (id int64) {
	if message == nil || message.Content == "" || message.Author.Id <= 0 {
		return ErrWrongData
	}
	row := s.db.QueryRow(`
		INSERT INTO messages (AuthorId, CreatedAt, Content)
		VALUES (?,DATETIME('now','utc'),?)
		RETURNING Id`,
		message.Author.Id,
		message.Content)
	if row.Scan(&id) != nil {
		return ErrMessageNotCreated
	}
	return
}
