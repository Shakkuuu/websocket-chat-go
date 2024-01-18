package entity

import "golang.org/x/net/websocket"

type ChatRoom struct {
	ID      string
	Clients map[*websocket.Conn]bool
}

type Data struct {
	Rooms   []string
	RoomID  string
	Message string
}

type Message struct {
	RoomID  string `json:"roomID"`
	Message string `json:"message"`
	Name    string `json:"name"`
}
