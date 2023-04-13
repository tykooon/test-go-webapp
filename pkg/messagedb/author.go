package messagedb

import (
	"time"
)

type Author struct {
	Id        int64     //`json:"-"`
	Name      string    `json:"name:"`
	CreatedAt time.Time `json:"created"`
}

func NewAuthor(name string) *Author {
	return &Author{
		Name:      name,
		CreatedAt: time.Now(),
	}
}

func (a Author) GetId() int64 {
	return a.Id
}
