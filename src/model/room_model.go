package model

import (
	"websocket-chat/db"
	"websocket-chat/entity"

	"golang.org/x/net/websocket"
)

var rooms = make(map[string]*entity.ChatRoom) // 作成された各ルームを格納

// DBから作成済みのRoomIDを持ってきて、作成する。
func RoomInit() error {
	db := db.GetDB()
	var r []entity.DBRoom

	err := db.Find(&r).Error
	if err != nil {
		return err
	}

	for _, room := range r {
		CreateRoom(room.RoomID)
	}

	return nil
}

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
