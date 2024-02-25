package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"websocket-chat/entity"
	"websocket-chat/model"

	"golang.org/x/net/websocket"
	"gorm.io/gorm"
)

var sentmessage = make(chan entity.Message) // 各クライアントに送信するためのメッセージのチャネル
var chatlogfile *os.File

// mainでOpenしたログファイルを変数に入れる
func ChatLogInit(f *os.File) {
	chatlogfile = f
}

// "YYYY-MM-DD HH-MM-SS"に変換
func timeToStr(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

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
			log.Printf("SessionToGetName error: %v\n", err)
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
		// セッションのユーザー取得
		user, err = model.GetUserByName(un)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("model.GetUserByName error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "ユーザーが見つかりませんでした。"

			err = troomtop.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}
		if err != nil {
			log.Printf("model.GetUserByName error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "データベースとの接続に失敗しました。"

			err = troomtop.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// Room作成
		err = model.DBCreateRoom(roomid)
		if err != nil {
			log.Printf("model.DBCreateRoom error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "データベースとの接続に失敗しました。"

			err = troomtop.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}
		room := model.CreateRoom(roomid)

		// 参加中のルーム一覧にMasterとして追加
		var proom entity.ParticipatingRoom
		proom = entity.ParticipatingRoom{
			RoomID:   room.ID,
			IsMaster: true,
			UserName: user.Name,
		}
		err = model.AddParticipatingRoom(&proom)
		if err != nil {
			log.Printf("model.AddParticipatingRoom error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "データベースとの接続に失敗しました。"

			err = troomtop.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// メッセージをテンプレートに渡す
		var data entity.Data
		data.Message = "ルーム " + roomid + " が作成されました。"

		err = troomtop.Execute(w, data)
		if err != nil {
			log.Printf("Excute error:%v\n", err)
			http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
			return
		}
	case http.MethodHead:
		fmt.Fprintln(w, "Thank you monitor.")
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

		_, exists := rooms[roomid]
		if !exists { // 指定した部屋が存在していなかったら
			log.Println("This room was not found")

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

		err = troom.Execute(w, nil)
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
		_, exists := rooms[roomid]
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
			log.Printf("SessionToGetName error: %v\n", err)
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
		// セッションのユーザー取得
		user, err = model.GetUserByName(un)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("model.GetUserByName error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "ユーザーが見つかりませんでした。"

			err = troomtop.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}
		if err != nil {
			log.Printf("model.GetUserByName error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "データベースとの接続に失敗しました。"

			err = troomtop.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		var proom entity.ParticipatingRoom
		prooms, err := model.GetParticipatingRoomByName(user.Name)
		if err != nil {
			log.Printf("model.GetParticipatingRoomByName error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "データベースとの接続に失敗しました。"

			err = troomtop.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		for _, proom = range prooms {
			if proom.RoomID == roomid {
				break
			}
		}

		if !proom.IsMaster {
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
		err = model.DeleteParticipatingRoomByRoomID(roomid)
		if err != nil {
			log.Printf("model.DeleteParticipatingRoomByRoomID error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "データベースとの接続に失敗しました。"

			err = troomtop.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// 部屋削除
		err = model.DeleteRoom(roomid)
		if err != nil {
			log.Printf("model.DeleteRoom error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "データベースとの接続に失敗しました。"

			err = troomtop.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

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
		log.Printf("This room was not found\n")
		return
	}

	var user entity.User
	// ユーザー取得
	user, err = model.GetUserByName(msg.Name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("model.GetUserByName error: %v\n", err)
		log.Printf("User Not Found: %v\n", err)
		return
	}
	if err != nil {
		log.Printf("model.GetUserByName error: %v\n", err)
		log.Printf("GetUserByName error: %v\n", err)
		return
	}

	var check bool = false
	// 既に参加しているかどうかを確認
	var proom entity.ParticipatingRoom
	prooms, err := model.GetParticipatingRoomByName(user.Name)
	if err != nil {
		log.Printf("GetParticipatingRoomByName error: %v\n", err)
		return
	}

	for _, proom = range prooms {
		if proom.RoomID == room.ID {
			check = true
		}
	}

	if !check {
		// 参加中のルーム一覧に参加者として追加
		var proom entity.ParticipatingRoom = entity.ParticipatingRoom{
			RoomID:   room.ID,
			IsMaster: false,
			UserName: user.Name,
		}
		err = model.AddParticipatingRoom(&proom)
		if err != nil {
			log.Printf("AddParticipatingRoom error: %v\n", err)
			return
		}
	}

	// Roomに参加
	room.Clients[ws] = msg.Name
	fmt.Println(room.Clients) // 参加者一覧 デバッグ用

	// 参加しているユーザー一覧とオンラインのユーザー一覧の取得
	allusersChan := make(chan interface{})
	onlineusersChan := make(chan interface{})
	var allusers []string
	var onlineusers []string
	go func() {
		aus, err := model.GetAllUsers(msg.RoomID)
		if err != nil {
			err = fmt.Errorf("GetAllUsers error: %v", err)
			allusersChan <- err
			return
		}
		allusersChan <- aus
	}()
	go func() {
		ous, err := model.GetOnlineUsers(msg.RoomID)
		if err != nil {
			err = fmt.Errorf("GetOnlineUsers error: %v", err)
			allusersChan <- err
			return
		}
		onlineusersChan <- ous
	}()

	v1 := <-allusersChan
	v2 := <-onlineusersChan
	switch t1 := v1.(type) {
	case error:
		log.Println(t1)
		return
	case []string:
		allusers = t1
	}

	switch t2 := v2.(type) {
	case error:
		log.Println(t2)
		return
	case []string:
		onlineusers = t2
	}

	// Roomに参加したことをそのRoomのクライアントにブロードキャスト
	entermsg := entity.Message{RoomID: room.ID, Message: msg.Name + "が入室しました", Name: "Server", ToName: "", AllUsers: allusers, OnlineUsers: onlineusers}
	sentmessage <- entermsg

	// サーバ側からクライアントにWellcomeメッセージを送信
	err = websocket.JSON.Send(ws, entity.Message{RoomID: room.ID, Message: "サーバ" + room.ID + "へようこそ", Name: "Server", ToName: msg.Name, AllUsers: nil, OnlineUsers: nil})
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

				// 参加しているユーザー一覧とオンラインのユーザー一覧の取得
				allusers, err := model.GetAllUsers(msg.RoomID)
				if err != nil {
					log.Printf("GetAllUsers error: %v\n", err)
					return
				}
				onlineusers, err := model.GetOnlineUsers(msg.RoomID)
				if err != nil {
					log.Printf("GetOnlineUsers error: %v\n", err)
					return
				}

				// そのクライアントがRoomから退出したことをそのRoomにブロードキャスト
				exitmsg := entity.Message{RoomID: msg.RoomID, Message: msg.Name + "が退出しました", Name: "Server", ToName: "", AllUsers: allusers, OnlineUsers: onlineusers}
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

		// チャットログを出力と保存 日時、サーバー名、ユーザー名、宛先、メッセージ
		replaceNlMsg := strings.ReplaceAll(msg.Message, "\n", " ") // 改行があるとログが改行されてしまうため、改行を削除
		chatlog := fmt.Sprintf("%s: [S%s] From(%s) To (%s) Msg(%s)\n", timeToStr(time.Now()), msg.RoomID, msg.Name, msg.ToName, replaceNlMsg)
		fmt.Print(chatlog)
		fmt.Fprint(chatlogfile, chatlog)

		if msg.ToName != "" {
			// 接続中のクライアントにメッセージを送る
			for client, name := range room.Clients {
				if msg.ToName == name || msg.Name == name {
					// メッセージを返信する
					err := websocket.JSON.Send(client, entity.Message{RoomID: room.ID, Message: msg.Message, Name: msg.Name, ToName: msg.ToName, AllUsers: msg.AllUsers, OnlineUsers: msg.OnlineUsers})
					if err != nil {
						log.Printf("Send error:%v\n", err)
					}
				}
			}
		} else {
			// 接続中のクライアントにメッセージを送る
			for client := range room.Clients {
				// メッセージを返信する
				err := websocket.JSON.Send(client, entity.Message{RoomID: room.ID, Message: msg.Message, Name: msg.Name, ToName: "", AllUsers: msg.AllUsers, OnlineUsers: msg.OnlineUsers})
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
			log.Printf("json.Marshal error: %v\n", err)
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
