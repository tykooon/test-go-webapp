package messagedb

import (
	"time"
)

type Message struct {
	Id        int64     //`json:"-"`
	Author    Author    `json:"author,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Content   string    `json:"content,omitempty"`
}

func NewMessage(author *Author, content string) *Message {
	return &Message{
		Author:    *author,
		CreatedAt: time.Now(),
		Content:   content,
	}
}

func (m *Message) GetId() int64 {
	return m.Id
}
