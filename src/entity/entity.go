package entity

import "golang.org/x/net/websocket"

// クライアントが参加するチャットルーム
type ChatRoom struct {
	ID      string
	Clients map[*websocket.Conn]string
}

// HTMLテンプレートに渡すためのデータ
type Data struct {
	Rooms   []string
	Users   []string
	RoomID  string
	Message string
}

// クライアントサーバ間でやりとりするメッセージ
type Message struct {
	RoomID  string `json:"roomID"`
	Message string `json:"message"`
	Name    string `json:"name"`
	ToName  string `json:"toname"`
}

// Room内のユーザー一覧送信用
type SentRoomUsersList struct {
	UsersList []string `json:"userslist"`
}

// ルーム一覧送信用
type SentRoomsList struct {
	RoomsList []string `json:"roomslist"`
}

// ユーザー名送信用
type SentUser struct {
	Name string `json:"name"`
}

// ユーザー管理
type User struct {
	Name               string
	Password           string
	ParticipatingRooms map[*ChatRoom]bool
}
