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
	ID       int    `gorm:"unique"`
	Name     string `gorm:"unique"`
	Password string
}

type ParticipatingRoom struct {
	ID       int `gorm:"unique"`
	RoomID   string
	IsMaster bool
	UserName string
}

// type User struct {
// 	ID                 int    `gorm:"unique"`
// 	Name               string `gorm:"unique"`
// 	Password           string
// 	ParticipatingRooms map[*ChatRoom]bool `gorm:"type:JSON"`
// 	CreatedAt          time.Time
// }

// 参加中の部屋情報
// type ParticipatingRoom struct {
// 	ID int `gorm:"unique"`
// 	Room ChatRoom
// 	IsMaster bool
// 	UserName string
// }
