package model

import (
	"websocket-chat/entity"

	"golang.org/x/net/websocket"
)

// Room作成
func CreateRoom(roomid string, rooms map[string]*entity.ChatRoom) *entity.ChatRoom {
	room := &entity.ChatRoom{
		ID:      roomid,
		Clients: make(map[*websocket.Conn]string),
	}
	rooms[roomid] = room
	return room
}
