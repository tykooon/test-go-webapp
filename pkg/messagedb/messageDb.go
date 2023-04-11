package messagedb

import (
	"database/sql"
)

const (
	ErrWrongData int64 = -(iota + 1)
	ErrEntryAlreadyExists
	ErrEntryNotFound
	ErrAuthorNotCreated
	ErrMessageNotCreated
)

type MessageDB struct {
	db                *sql.DB
	authorByNameQuery string // Args: (name string) Returns: (authorId int64)
	authorById        string // Args: (id int64) Returns: (name string, createdAt string)
	messageById       string // Args: (id int64) Returns: (authorId int64, createdAt string, content string)
	allMessages       string // Args: <empty> Returns: [](authorId int64, createdAt string, content string)
	insertAuthor      string // Args: (name string) Returns: (authorId int64)
	insertMessage     string // Args: (authorId int64, content string) Returns: (messageId int64)
}

func NewMessageDB(db *sql.DB, dict map[string]string) *MessageDB {
	return &MessageDB{
		db:                db,
		authorByNameQuery: dict["author_by_name"],
		authorById:        dict["author_by_id"],
		messageById:       dict["message_by_id"],
		allMessages:       dict["all_messages"],
		insertAuthor:      dict["insert_author"],
		insertMessage:     dict["insert_message"],
	}
}

func (m *MessageDB) GetAuthorByName(name string) *Author {
	author := NewAuthor(name)
	row := m.db.QueryRow(m.authorByNameQuery,
		name)
	if row.Scan(&author.id) != nil {
		return nil
	}
	return author
}

func (m *MessageDB) GetAuthorById(id int64) *Author {
	author := &Author{id: id}
	row := m.db.QueryRow(m.authorById,
		id)
	if row.Scan(&author.Name, &author.CreatedAt) != nil {
		return nil
	}
	return author
}

func (m *MessageDB) GetMessageById(id int64) *Message {
	message := &Message{id: id}
	var authId int64
	row := m.db.QueryRow(m.messageById,
		id)
	if err := row.Scan(&authId, &message.CreatedAt, &message.Content); err != nil {
		return nil
	}
	message.Author = *m.GetAuthorById(authId)
	return message
}

func (m *MessageDB) GetAllMessages() []*Message {
	rows, err := m.db.Query(m.allMessages)
	if err != nil {
		return nil
	}
	defer rows.Close()
	messages := make([]*Message, 0)
	for rows.Next() {
		message := &Message{}
		if err := rows.Scan(&message.Author.id, &message.CreatedAt, &message.Content); err != nil {
			return nil
		}
		message.Author = *m.GetAuthorById(message.Author.id)
		messages = append(messages, message)
	}
	return messages
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
	row := m.db.QueryRow(m.insertAuthor,
		author.Name)
	if row.Scan(&id) != nil {
		return ErrAuthorNotCreated
	}
	author.id = id
	return
}

func (m *MessageDB) InsertMessage(message *Message) (id int64) {
	if message == nil || message.Content == "" || message.Author.id <= 0 {
		return ErrWrongData
	}
	row := m.db.QueryRow(m.insertMessage,
		message.Author.id,
		message.Content)
	if row.Scan(&id) != nil {
		return ErrMessageNotCreated
	}
	return
}
