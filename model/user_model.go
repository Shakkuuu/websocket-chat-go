package model

import "websocket-chat/entity"

// TODO:dbにしてみる？
var users = []entity.User{
	{Name: "匿名", Password: "qawsedrftgyhujikolp"},
}

// ユーザー一覧取得
func GetUsers() []entity.User {
	return users
}

// ユーザー追加
func AddUser(user entity.User) {
	users = append(users, user)
}
