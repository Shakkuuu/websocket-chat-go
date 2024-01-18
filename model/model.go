package model

import (
	"fmt"

	"websocket-chat/entity"

	"golang.org/x/net/websocket"
)

func CreateRoom(roomid string, rooms map[string]*entity.ChatRoom) {
	room := &entity.ChatRoom{
		ID:      roomid,
		Clients: make(map[*websocket.Conn]bool),
	}
	fmt.Printf("room %v が作成されました\n", room.ID)
	rooms[roomid] = room
}
