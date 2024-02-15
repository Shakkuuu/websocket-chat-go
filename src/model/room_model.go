package model

import (
	"websocket-chat/entity"

	"golang.org/x/net/websocket"
)

var rooms = make(map[string]*entity.ChatRoom) // 作成された各ルームを格納

// Room一覧取得
func GetRooms() map[string]*entity.ChatRoom {
	return rooms
}

// Room作成
func CreateRoom(roomid string) *entity.ChatRoom {
	room := &entity.ChatRoom{
		ID:      roomid,
		Clients: make(map[*websocket.Conn]string),
	}
	rooms[roomid] = room
	return room
}

// Room削除
func DeleteRoom(roomid string) {
	delete(rooms, roomid)
}
