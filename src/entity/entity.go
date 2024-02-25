package entity

import (
	"golang.org/x/net/websocket"
)

// クライアントが参加するチャットルーム
type ChatRoom struct {
	ID      string
	Clients map[*websocket.Conn]string
}

// HTMLテンプレートに渡すためのデータ
type Data struct {
	Rooms   []string
	RoomID  string
	Name    string
	Message string
}

// クライアントサーバ間でやりとりするメッセージ
type Message struct {
	RoomID      string   `json:"roomID"`
	Message     string   `json:"message"`
	Name        string   `json:"name"`
	ToName      string   `json:"toname"`
	AllUsers    []string `json:"allusers"`
	OnlineUsers []string `json:"onlineusers"`
}

// // Room内のユーザー一覧送信用
// type SentRoomUsersList struct {
// 	UsersList []string `json:"userslist"`
// }

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
	ID       int    `gorm:"unique"`
	Name     string `gorm:"unique"`
	Password string
}

// ユーザーの参加中Room情報
type ParticipatingRoom struct {
	ID       int `gorm:"unique"`
	RoomID   string
	IsMaster bool
	UserName string
}

// DB保存用Room
type DBRoom struct {
	ID     int    `gorm:"unique"`
	RoomID string `gorm:"unique"`
}
