package entity

import "golang.org/x/net/websocket"

// クライアントが参加するチャットルーム
type ChatRoom struct {
	ID      string
	Clients map[*websocket.Conn]bool
}

// HTMLテンプレートに渡すためのデータ
type Data struct {
	Rooms   []string
	RoomID  string
	Message string
}

// クライアントサーバ間でやりとりするメッセージ
type Message struct {
	RoomID  string `json:"roomID"`
	Message string `json:"message"`
	Name    string `json:"name"`
}
