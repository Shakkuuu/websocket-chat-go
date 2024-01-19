package controller

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"websocket-chat/entity"
	"websocket-chat/model"

	"golang.org/x/net/websocket"
)

var rooms = make(map[string]*entity.ChatRoom) // 作成された各ルームを格納

var broadcast = make(chan entity.Message) // 各クライアントにブロードキャストするためのメッセージのチャネル

// indexページの表示
func Index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t, err := template.ParseFiles("view/index.html")
		if err != nil {
			log.Printf("template.ParseFiles error:%v\n", err)
			http.Error(w, "ページの読み込みに失敗しました。", http.StatusInternalServerError)
			return
		}

		// Room一覧をテンプレートに渡す
		var data entity.Data
		for k := range rooms {
			data.Rooms = append(data.Rooms, k)
		}

		err = t.Execute(w, data)
		if err != nil {
			log.Printf("Excute error:%v\n", err)
			http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		// POSTされた作成するRoomidをFormから受け取り
		r.ParseForm()
		roomid := r.FormValue("create_roomid")

		_, exists := rooms[roomid]
		if exists { // roomが既に存在していたら
			// 作成失敗メッセージ表示
			t, err := template.ParseFiles("view/index.html")
			if err != nil {
				log.Printf("template.ParseFiles error:%v\n", err)
				http.Error(w, "ページの読み込みに失敗しました。", http.StatusInternalServerError)
				return
			}

			// Room一覧とメッセージをテンプレートに渡す
			var data entity.Data
			for k := range rooms {
				data.Rooms = append(data.Rooms, k)
			}
			data.Message = "ルーム " + roomid + " は既にあります"

			err = t.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// Room作成
		model.CreateRoom(roomid, rooms)

		t, err := template.ParseFiles("view/index.html")
		if err != nil {
			log.Printf("template.ParseFiles error:%v\n", err)
			http.Error(w, "ページの読み込みに失敗しました。", http.StatusInternalServerError)
			return
		}

		// Room一覧とメッセージをテンプレートに渡す
		var data entity.Data
		for k := range rooms {
			data.Rooms = append(data.Rooms, k)
		}
		data.Message = "ルーム " + roomid + " が作成されました。"

		err = t.Execute(w, data)
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
		// POSTされたRoomIDをFormから取得
		r.ParseForm()
		roomid := r.URL.Query().Get("roomid")

		_, exists := rooms[roomid]
		if !exists { // 指定した部屋が存在していなかったら
			log.Printf("This room was not found")

			t, err := template.ParseFiles("view/index.html")
			if err != nil {
				log.Printf("template.ParseFiles error:%v\n", err)
				http.Error(w, "ページの読み込みに失敗しました。", http.StatusInternalServerError)
				return
			}

			// Roomの一覧とメッセージをテンプレートに渡す
			var data entity.Data
			for k := range rooms {
				data.Rooms = append(data.Rooms, k)
			}
			data.Message = "そのIDのルームは見つかりませんでした。"

			err = t.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		t, err := template.ParseFiles("view/room.html")
		if err != nil {
			log.Printf("template.ParseFiles error:%v\n", err)
			http.Error(w, "ページの読み込みに失敗しました。", http.StatusInternalServerError)
			return
		}

		err = t.Execute(w, nil)
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

	// 部屋が生きているかどうか(なくてもいいかも)
	room, exists := rooms[msg.RoomID]
	if !exists {
		log.Printf("This room was not found")
		return
	}

	// Roomに参加
	room.Clients[ws] = true
	fmt.Println(room.Clients) // 参加者一覧 デバッグ用

	// Roomに参加したことをそのRoomのクライアントにブロードキャスト
	entermsg := entity.Message{RoomID: room.ID, Message: msg.Name + "が入室しました", Name: "Server"}
	broadcast <- entermsg

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
				broadcast <- exitmsg
				break
			}
			log.Printf("Receive error:%v\n", err)
		}

		// goroutineでチャネルを待っているとこへメッセージを渡す
		broadcast <- msg
	}
}

// goroutineでメッセージのチャネルが来るまで待ち、Roomにブロードキャストする
func HandleMessages() {
	for {
		// broadcastチャネルからメッセージを受け取る
		msg := <-broadcast
		room, exists := rooms[msg.RoomID]
		if !exists {
			continue
		}
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
