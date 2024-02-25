package model

import (
	"fmt"
	"websocket-chat/db"
	"websocket-chat/entity"

	"golang.org/x/net/websocket"
)

var rooms = make(map[string]*entity.ChatRoom) // 作成された各ルームを格納
var err error

// DBから作成済みのRoomIDを持ってきて、作成する。
func RoomInit() error {
	db := db.GetDB()
	var r []entity.DBRoom

	err = db.Find(&r).Error
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

// Room作成(db)
func DBCreateRoom(roomid string) error {
	db := db.GetDB()
	var r entity.DBRoom
	r.RoomID = roomid

	err = db.Create(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// Room削除
func DeleteRoom(roomid string) error {
	db := db.GetDB()

	var r entity.DBRoom

	err = db.Where("room_id = ?", roomid).Delete(&r).Error
	if err != nil {
		return err
	}

	delete(rooms, roomid)

	return nil
}

// 参加しているユーザー一覧の取得
func GetAllUsers(roomid string) ([]string, error) {
	var allusers []string

	// Roomに参加しているユーザーの取得
	allusers = append(allusers, "匿名")
	allprooms, err := GetParticipatingRoomByRoomID(roomid)
	if err != nil {
		err = fmt.Errorf("GetParticipatingRoomByRoomID error: %v", err)
		return allusers, err
	}

	// ユーザーを格納
	for _, prm := range allprooms {
		allusers = append(allusers, prm.UserName)
	}

	return allusers, err

}

// オンラインのユーザー一覧の取得
func GetOnlineUsers(roomid string) ([]string, error) {
	var onlineusers []string

	// オンラインのユーザー取得
	onlineusers = append(onlineusers, "匿名")

	// Room一覧取得
	rooms = GetRooms()

	// roomがあるか再度確認
	room, exists := rooms[roomid]
	if !exists {
		err = fmt.Errorf("this room was not found")
		return onlineusers, err
	}

	// Room内のユーザーを格納
	for _, user := range room.Clients {
		onlineusers = append(onlineusers, user)
	}

	return onlineusers, nil
}
