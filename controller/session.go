package controller

import (
	"github.com/gorilla/sessions"
)

const SESSION_NAME string = "Shakkuuu-websocket-chat-go"

var session *sessions.Session
var store *sessions.CookieStore
var err error

func SessionInit(sessionKey string) {
	store = sessions.NewCookieStore([]byte(sessionKey))
	session = sessions.NewSession(store, SESSION_NAME)
}
