package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

const SESSION_NAME string = "Shakkuuu-websocket-chat-go"

var session *sessions.Session
var store *sessions.CookieStore
var err error

// セッションの初期化
func SessionInit(sessionKey string) {
	store = sessions.NewCookieStore([]byte(sessionKey))
	session = sessions.NewSession(store, SESSION_NAME)
}

// セッションからユーザー名を取得
func SessionToGetName(r *http.Request) (string, error) {
	// セッション読み取り
	session, err = store.Get(r, SESSION_NAME)
	if err != nil {
		log.Printf("store.Get error: %v", err)
		return "", err
	}

	username := session.Values["username"]
	if username == nil {
		fmt.Println("セッションなし")
		err = fmt.Errorf("セッションなし")
		return "", err
	}

	un := username.(string)
	return un, nil
}
