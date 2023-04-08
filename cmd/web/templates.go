package main

import (
	"html/template"
	"time"
)

type templates struct {
	home *template.Template
}

func initTemplates() *templates {
	res := templates{}
	res.home = template.New("home")
	fmap := template.FuncMap(map[string]any{
		"timeformat": timeformat,
	})
	res.home.Funcs(fmap)
	res.home = template.Must(res.home.ParseFiles("template/base.html", "template/home.html"))
	return &res
}

func timeformat(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("02 Jan 2006 15:04")
}
