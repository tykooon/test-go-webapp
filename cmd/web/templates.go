package main

import (
	"html/template"
)

type templates struct {
	home *template.Template
}

func initTemplates() *templates {
	res := templates{}
	res.home = template.Must(template.New("home").ParseFiles("template/base.html", "template/home.html"))
	return &res
}
