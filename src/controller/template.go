package controller

import (
	"html/template"
	"log"
)

var tlogin *template.Template
var troom *template.Template
var troomtop *template.Template
var tsignup *template.Template
var tusermenu *template.Template

func TemplateInit() {
	tlogin, err = template.ParseFiles("view/login.html")
	if err != nil {
		log.Printf("template.ParseFiles error:%v\n", err)
		return
	}
	troom, err = template.ParseFiles("view/room.html")
	if err != nil {
		log.Printf("template.ParseFiles error:%v\n", err)
		return
	}
	troomtop, err = template.ParseFiles("view/roomtop.html")
	if err != nil {
		log.Printf("template.ParseFiles error:%v\n", err)
		return
	}
	tsignup, err = template.ParseFiles("view/signup.html")
	if err != nil {
		log.Printf("template.ParseFiles error:%v\n", err)
		return
	}
	tusermenu, err = template.ParseFiles("view/usermenu.html")
	if err != nil {
		log.Printf("template.ParseFiles error:%v\n", err)
		return
	}
}
