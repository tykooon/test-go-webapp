package webappmodel

import (
	"html/template"
	"time"
)

type templates struct {
	home     *template.Template
	messages *template.Template
	send     *template.Template
}

func initTemplates() *templates {
	res := templates{}

	res.home = template.New("home")
	res.messages = template.New("messages")
	res.send = template.New("send")

	fmap := template.FuncMap(map[string]any{
		"timeformat": timeformat,
	})

	res.home = template.Must(res.home.ParseFiles("template/base.html", "template/home.html"))

	res.messages.Funcs(fmap)
	res.messages = template.Must(res.messages.ParseFiles("template/base.html", "template/messages.html"))

	res.send = template.Must(res.send.ParseFiles("template/base.html", "template/send.html"))

	return &res
}

func timeformat(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("02 Jan 2006 15:04")
}
