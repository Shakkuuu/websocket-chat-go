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

// 名前からのユーザー取得処理
func GetUserByName(username string) (entity.User, error) {
	db := db.GetDB()
	var u entity.User

	err := db.Where("name = ?", username).First(&u).Error
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

// Nameからのユーザーデータ更新処理
// func PutUserByName(u *entity.User, username string) error {
// 	db := db.GetDB()

// 	err := db.Where("username = ?", username).Model(&u).Updates(&u).Error
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// Nameからのユーザー削除処理
func DeleteUser(username string) error {
	db := db.GetDB()

	var u entity.User

	err := db.Where("name = ?", username).Delete(&u).Error
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
func DeleteParticipatingRoomByRoomID(roomid string) error {
	db := db.GetDB()

	var p entity.ParticipatingRoom

	err := db.Where("room_id = ?", roomid).Delete(&p).Error
	if err != nil {
		return err
	}

	return nil
}

// 参加中RoomをUserNameから削除
func DeleteParticipatingRoomByUserName(username string) error {
	db := db.GetDB()

	var p entity.ParticipatingRoom

	err := db.Where("user_name = ?", username).Delete(&p).Error
	if err != nil {
		return err
	}

	return nil
}
