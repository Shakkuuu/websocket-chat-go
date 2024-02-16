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

func TemplateInit() error {
	tlogin, err = template.ParseFiles("view/login.html")
	if err != nil {
		log.Printf("template.ParseFiles error:%v\n", err)
		return err
	}
	troom, err = template.ParseFiles("view/room.html")
	if err != nil {
		log.Printf("template.ParseFiles error:%v\n", err)
		return err
	}
	troomtop, err = template.ParseFiles("view/roomtop.html")
	if err != nil {
		log.Printf("template.ParseFiles error:%v\n", err)
		return err
	}
	tsignup, err = template.ParseFiles("view/signup.html")
	if err != nil {
		log.Printf("template.ParseFiles error:%v\n", err)
		return err
	}
	tusermenu, err = template.ParseFiles("view/usermenu.html")
	if err != nil {
		log.Printf("template.ParseFiles error:%v\n", err)
		return err
	}
	return nil
}
