package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"websocket-chat/entity"
	"websocket-chat/model"

	"golang.org/x/net/websocket"
)

var sentmessage = make(chan entity.Message) // 各クライアントに送信するためのメッセージのチャネル

// roomtopページの表示
func RoomTop(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err = troomtop.Execute(w, nil)
		if err != nil {
			log.Printf("Excute error:%v\n", err)
			http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		// POSTされた作成するRoomidをFormから受け取り
		r.ParseForm()
		roomid := r.FormValue("create_roomid")

		// Room一覧取得
		rooms := model.GetRooms()

		_, exists := rooms[roomid]
		if exists { // roomが既に存在していたら
			// 作成失敗メッセージ表示
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "ルーム " + roomid + " は既にあります"

			err = troomtop.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// セッション読み取り
		un, err := SessionToGetName(r)
		if err != nil {
			log.Printf("SessionToGetName error: %v", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "再ログインしてください"

			err = tlogin.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		var user entity.User
		var check bool
		// ユーザー一覧取得
		users := model.GetUsers()
		// ユーザーリストからセッションと一致するユーザーを持ってくる
		for _, v := range users {
			if un == v.Name {
				user = v
				check = true
			}
		}

		if !check {
			fmt.Println("セッション問題あり")
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "再ログインしてください"

			err = tlogin.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// Room作成
		room := model.CreateRoom(roomid)

		// 参加中のルーム一覧にMasterとして追加
		user.ParticipatingRooms[room] = true

		// メッセージをテンプレートに渡す
		var data entity.Data
		data.Message = "ルーム " + roomid + " が作成されました。"

		err = troomtop.Execute(w, data)
		if err != nil {
			log.Printf("Excute error:%v\n", err)
			http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
			return
		}
	default:
		fmt.Fprintln(w, "Method not allowed")
		http.Error(w, "そのメソッドは許可されていません。", http.StatusMethodNotAllowed)
		return
	}
}

// Room内のページ
func Room(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// クエリ読み取り
		r.ParseForm()
		roomid := r.URL.Query().Get("roomid")

		// Room一覧取得
		rooms := model.GetRooms()

		room, exists := rooms[roomid]
		if !exists { // 指定した部屋が存在していなかったら
			log.Printf("This room was not found")

			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "そのIDのルームは見つかりませんでした。"

			err = troomtop.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// 参加しているユーザー一覧をテンプレートに渡す
		var data entity.Data
		for _, username := range room.Clients {
			data.Users = append(data.Users, username)
		}

		err = troom.Execute(w, data)
		if err != nil {
			log.Printf("Excute error:%v\n", err)
			http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
			return
		}
	default:
		fmt.Fprintln(w, "Method not allowed")
		http.Error(w, "そのメソッドは許可されていません。", http.StatusMethodNotAllowed)
		return
	}
}

// Room削除
func DeleteRoom(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// クエリ読み取り
		r.ParseForm()
		roomid := r.URL.Query().Get("roomid")

		// Room一覧取得
		rooms := model.GetRooms()

		// roomがあるか確認
		room, exists := rooms[roomid]
		if !exists {
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "そのIDのルームは見つかりませんでした。"

			err = troomtop.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// セッション読み取り
		un, err := SessionToGetName(r)
		if err != nil {
			log.Printf("SessionToGetName error: %v", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "再ログインしてください"

			err = tlogin.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		var user entity.User
		var check bool
		// ユーザー一覧取得
		users := model.GetUsers()
		// ユーザーリストからセッションと一致するユーザーを持ってくる
		for _, v := range users {
			if un == v.Name {
				user = v
				check = true
			}
		}

		if !check {
			fmt.Println("セッション問題あり")
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "再ログインしてください"

			err = tlogin.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		if !user.ParticipatingRooms[room] {
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "部屋の作成者ではないため、部屋を削除できませんでした。"

			err = troomtop.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// ユーザーの参加中ルームリストからも削除
		for _, u := range users {
			delete(u.ParticipatingRooms, room)
		}

		// 部屋削除
		model.DeleteRoom(roomid)

		// メッセージをテンプレートに渡す
		var data entity.Data
		data.Message = "部屋を削除しました。"

		err = troomtop.Execute(w, data)
		if err != nil {
			log.Printf("Excute error:%v\n", err)
			http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
			return
		}
	default:
		fmt.Fprintln(w, "Method not allowed")
		http.Error(w, "そのメソッドは許可されていません。", http.StatusMethodNotAllowed)
		return
	}
}

// WebsocketでRoom参加後のコネクション確立
func HandleConnection(ws *websocket.Conn) {
	defer ws.Close()

	// クライアントから参加する部屋が指定されたメッセージ受信
	var msg entity.Message
	err := websocket.JSON.Receive(ws, &msg)
	if err != nil {
		log.Printf("Receive room ID error:%v\n", err)
		return
	}

	// Room一覧取得
	rooms := model.GetRooms()

	// 部屋が存在しているかどうか(なくてもいいかも)
	room, exists := rooms[msg.RoomID]
	if !exists {
		log.Printf("This room was not found")
		return
	}

	var user entity.User
	// ユーザー一覧取得
	users := model.GetUsers()
	// ユーザーリストからメッセージのNameと一致するユーザーを持ってくる
	for _, u := range users {
		if msg.Name == u.Name {
			user = u
		}
	}

	var check bool
	// 既に参加しているかどうかを確認
	for r := range user.ParticipatingRooms {
		if r == room {
			check = true
		}
	}
	if !check {
		// 参加中のルーム一覧に参加者として追加
		user.ParticipatingRooms[room] = false
	}

	// Roomに参加
	room.Clients[ws] = msg.Name
	fmt.Println(room.Clients) // 参加者一覧 デバッグ用

	// Roomに参加したことをそのRoomのクライアントにブロードキャスト
	entermsg := entity.Message{RoomID: room.ID, Message: msg.Name + "が入室しました", Name: "Server"}
	sentmessage <- entermsg

	// サーバ側からクライアントにWellcomeメッセージを送信
	err = websocket.JSON.Send(ws, entity.Message{RoomID: room.ID, Message: "サーバ" + room.ID + "へようこそ", Name: "Server"})
	if err != nil {
		log.Printf("server wellcome Send error:%v\n", err)
	}

	// クライアントからメッセージが来るまで受信待ちする
	for {
		// クライアントからのメッセージを受信
		err = websocket.JSON.Receive(ws, &msg)
		if err != nil {
			if err.Error() == "EOF" { // Roomを退出したことを示すメッセージが来たら
				log.Printf("EOF error:%v\n", err)
				delete(room.Clients, ws) // Roomからそのクライアントを削除
				// そのクライアントがRoomから退出したことをそのRoomにブロードキャスト
				exitmsg := entity.Message{RoomID: msg.RoomID, Message: msg.Name + "が退出しました", Name: "Server"}
				sentmessage <- exitmsg
				break
			}
			log.Printf("Receive error:%v\n", err)
		}

		// goroutineでチャネルを待っているとこへメッセージを渡す
		sentmessage <- msg
	}
}

// goroutineでメッセージのチャネルが来るまで待ち、Roomにメッセージを送信する
func HandleMessages() {
	for {
		// sentmessageチャネルからメッセージを受け取る
		msg := <-sentmessage

		// Room一覧取得
		rooms := model.GetRooms()

		// 部屋が存在しているかどうか
		room, exists := rooms[msg.RoomID]
		if !exists {
			continue
		}

		if msg.ToName != "" {
			// 接続中のクライアントにメッセージを送る
			for client, name := range room.Clients {
				if msg.ToName == name || msg.Name == name {
					// メッセージを返信する
					err := websocket.JSON.Send(client, entity.Message{RoomID: room.ID, Message: msg.Message, Name: msg.Name, ToName: msg.ToName})
					if err != nil {
						log.Printf("Send error:%v\n", err)
					}
				}
			}
		} else {
			// 接続中のクライアントにメッセージを送る
			for client := range room.Clients {
				// メッセージを返信する
				err := websocket.JSON.Send(client, entity.Message{RoomID: room.ID, Message: msg.Message, Name: msg.Name})
				if err != nil {
					log.Printf("Send error:%v\n", err)
				}
			}
		}
	}
}

// Roomの一覧を返す
func RoomsList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var roomslist entity.SentRoomsList

		// Room一覧取得
		rooms := model.GetRooms()

		// Roomを格納
		for room := range rooms {
			roomslist.RoomsList = append(roomslist.RoomsList, room)
		}

		// jsonに変換
		sentjson, err := json.Marshal(roomslist)
		if err != nil {
			log.Printf("json.Marshal error: %v", err)
			http.Error(w, "json.Marshal error", http.StatusInternalServerError)
			return
		}

		// jsonで送信
		w.Header().Set("Content-Type", "application/json")
		w.Write(sentjson)

	default:
		fmt.Fprintln(w, "Method not allowed")
		http.Error(w, "そのメソッドは許可されていません。", http.StatusMethodNotAllowed)
		return
	}
}

// 参加中のRoomの一覧を返す
func JoinRoomsList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var joinroomslist entity.SentRoomsList

		// セッション読み取り
		un, err := SessionToGetName(r)
		if err != nil {
			log.Printf("SessionToGetName error: %v", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "再ログインしてください"

			err = tlogin.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		var user entity.User
		var check bool
		// ユーザー一覧取得
		users := model.GetUsers()
		// ユーザーリストからセッションと一致するユーザーを持ってくる
		for _, v := range users {
			if un == v.Name {
				user = v
				check = true
			}
		}

		if !check {
			fmt.Println("セッション問題あり")
			http.Error(w, "session is posiible", http.StatusUnauthorized)
			return
		}

		// joinRoomを格納
		for room := range user.ParticipatingRooms {
			joinroomslist.RoomsList = append(joinroomslist.RoomsList, room.ID)
		}

		fmt.Println(joinroomslist.RoomsList)

		// jsonに変換
		sentjson, err := json.Marshal(joinroomslist)
		if err != nil {
			log.Printf("json.Marshal error: %v", err)
			http.Error(w, "json.Marshal error", http.StatusInternalServerError)
			return
		}

		// jsonで送信
		w.Header().Set("Content-Type", "application/json")
		w.Write(sentjson)

	default:
		fmt.Fprintln(w, "Method not allowed")
		http.Error(w, "そのメソッドは許可されていません。", http.StatusMethodNotAllowed)
		return
	}
}
