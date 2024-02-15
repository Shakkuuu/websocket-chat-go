package model

import (
	"websocket-chat/db"
	"websocket-chat/entity"
)

// ユーザー一覧取得
func GetUsers() ([]entity.User, error) {
	db := db.GetDB()
	var u []entity.User

	err := db.Find(&u).Error
	if err != nil {
		return u, err
	}

	return u, nil
}

// ユーザー追加
func AddUser(u *entity.User) error {
	db := db.GetDB()

	err := db.Create(&u).Error
	if err != nil {
		return err
	}

	return nil
}

// usernameから参加中Room取得
func GetParticipatingRoom(username string) ([]entity.ParticipatingRoom, error) {
	db := db.GetDB()
	var p []entity.ParticipatingRoom

	err := db.Where("user_name = ?", username).Find(&p).Error
	if err != nil {
		return p, err
	}

	return p, nil
}

// 参加中Room追加
func AddParticipatingRoom(p *entity.ParticipatingRoom) error {
	db := db.GetDB()

	err := db.Create(&p).Error
	if err != nil {
		return err
	}

	return nil
}

// 参加中RoomをRoomIDから削除
func DeleteParticipatingRoom(roomid string) error {
	db := db.GetDB()

	var p entity.ParticipatingRoom

	err := db.Where("room_id = ?", roomid).Delete(&p).Error
	if err != nil {
		return err
	}

	return nil
}
