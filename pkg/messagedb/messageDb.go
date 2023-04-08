package messagedb

import (
	"database/sql"
	"time"
)

const (
	ErrWrongData int64 = -(iota + 1)
	ErrEntryAlreadyExists
	ErrEntryNotFound
	ErrAuthorNotCreated
	ErrMessageNotCreated
)

type MessageDB struct {
	db *sql.DB
}

func NewMessageDB(db *sql.DB) *MessageDB {
	return &MessageDB{db}
}

func (m *MessageDB) GetAuthorByName(name string) *Author {
	author := NewAuthor(name)
	row := m.db.QueryRow(
		`SELECT Id 
		 FROM authors
		 WHERE name =?`,
		name)
	if row.Scan(&author.id) != nil {
		return nil
	}
	return author
}

func (m *MessageDB) GetAuthorById(id int64) *Author {
	author := &Author{id: id}
	row := m.db.QueryRow(
		`SELECT Name, CreatedAt 
         FROM authors
         WHERE id =?`,
		id)
	if row.Scan(&author.Name, &author.CreatedAt) == nil {
		return nil
	}
	return author
}

func (m *MessageDB) GetMessageById(id int64) *Message {
	message := &Message{id: id}
	timeStr := ""
	row := m.db.QueryRow(
		`SELECT AuthorId, CreatedAt, Content
		 FROM messages
		 WHERE id =?`,
		id)
	if err := row.Scan(&message.AuthorId, &timeStr, &message.Content); err != nil {
		//message.Content = err.Error()
		return nil //message
	}
	message.CreatedAt, _ = time.Parse(time.DateTime, timeStr[:19])
	return message
}

func (m *MessageDB) Insert(authorName string, content string) int64 {
	if authorName == "" || content == "" {
		return ErrWrongData
	}
	author := m.GetAuthorByName(authorName)
	if author == nil {
		author = NewAuthor(authorName)
		author.id = m.InsertAuthor(author)
		if author.id == ErrAuthorNotCreated {
			return ErrAuthorNotCreated
		}
	}
	message := NewMessage(author, content)
	messId := m.InsertMessage(message)
	if messId == ErrMessageNotCreated {
		return ErrMessageNotCreated
	}
	return messId
}

func (m *MessageDB) InsertAuthor(author *Author) (id int64) {
	if author.Name == "" {
		return ErrWrongData
	}
	if m.GetAuthorByName(author.Name) != nil {
		return ErrEntryAlreadyExists
	}
	row := m.db.QueryRow(
		`INSERT INTO authors (Name, CreatedAt) 
		 VALUES (?, ?)
		 RETURNING Id`,
		author.Name,
		time.Now())
	if row.Scan(&id) != nil {
		return ErrAuthorNotCreated
	}
	author.id = id
	return
}

func (m *MessageDB) InsertMessage(message *Message) (id int64) {
	if message.Content == "" || message.AuthorId <= 0 {
		return ErrWrongData
	}
	row := m.db.QueryRow(
		`INSERT INTO messages (AuthorId, CreatedAt, Content)  
         VALUES (?,?,?)
         RETURNING Id`,
		message.AuthorId,
		time.Now(),
		message.Content)
	if row.Scan(&id) != nil {
		return ErrMessageNotCreated
	}
	return
}

// func timeformat(t time.Time) string {
// 	if t.IsZero() {
// 		return ""
// 	}
// 	return t.Format("02 Jan 2006 15:04")
// }
