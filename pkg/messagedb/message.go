package messagedb

import (
	"time"
)

type Message struct {
	id        int64     //`json:"-"`
	AuthorId  int64     `json:"author_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Content   string    `json:"content,omitempty"`
}

func NewMessage(author *Author, content string) *Message {
	return &Message{
		AuthorId:  author.id,
		CreatedAt: time.Now(),
		Content:   content,
	}
}

func (m *Message) GetId() int64 {
	return m.id
}
