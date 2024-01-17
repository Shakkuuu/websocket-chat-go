package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

type ChatRoom struct {
	ID      string
	Clients map[*websocket.Conn]bool
}

type Data struct {
	Rooms   []string
	RoomID  string
	Message string
}

var rooms = make(map[string]*ChatRoom)

var broadcast = make(chan Message)

type Message struct {
	RoomID  string `json:"roomID"`
	Message string `json:"message"`
	Name    string `json:"name"`
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/room", room)
	http.Handle("/ws", websocket.Handler(handleConnection))
	go handleMessages()

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Printf("ListenAndServe error:%v\n", err)
		os.Exit(1)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t, err := template.ParseFiles("index.html")
		if err != nil {
			log.Printf("template.ParseFiles error:%v\n", err)
			http.Error(w, "ページの読み込みに失敗しました。", http.StatusInternalServerError)
			return
		}

		var data Data
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
		r.ParseForm()
		roomid := r.FormValue("create_roomid")

		_, exists := rooms[roomid]
		if exists { // roomが既に存在していたら
			// 作成失敗メッセージ表示
			t, err := template.ParseFiles("index.html")
			if err != nil {
				log.Printf("template.ParseFiles error:%v\n", err)
				http.Error(w, "ページの読み込みに失敗しました。", http.StatusInternalServerError)
				return
			}

			var data Data
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

		createRoom(roomid)

		t, err := template.ParseFiles("index.html")
		if err != nil {
			log.Printf("template.ParseFiles error:%v\n", err)
			http.Error(w, "ページの読み込みに失敗しました。", http.StatusInternalServerError)
			return
		}

		var data Data
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

func room(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		r.ParseForm()
		roomid := r.URL.Query().Get("roomid")

		_, exists := rooms[roomid]
		if !exists { // 指定した部屋が存在していなかったら
			log.Printf("This room was not found")

			t, err := template.ParseFiles("index.html")
			if err != nil {
				log.Printf("template.ParseFiles error:%v\n", err)
				http.Error(w, "ページの読み込みに失敗しました。", http.StatusInternalServerError)
				return
			}

			var data Data
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

		t, err := template.ParseFiles("room.html")
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

func createRoom(roomid string) {
	room := &ChatRoom{
		ID:      roomid,
		Clients: make(map[*websocket.Conn]bool),
	}
	fmt.Printf("room %v が作成されました\n", room.ID)
	rooms[roomid] = room
}

func handleConnection(ws *websocket.Conn) {
	defer ws.Close()

	var msg Message
	err := websocket.JSON.Receive(ws, &msg)
	if err != nil {
		log.Printf("Receive room ID error:%v\n", err)
		return
	}

	room, exists := rooms[msg.RoomID]
	if !exists {
		log.Printf("This room was not found")
		return
	}

	entermsg := Message{RoomID: room.ID, Message: msg.Name + "が入室しました", Name: "Server"}
	broadcast <- entermsg

	err = websocket.JSON.Send(ws, Message{RoomID: room.ID, Message: "サーバ" + room.ID + "へようこそ", Name: "Server"})
	if err != nil {
		log.Printf("server wellcome Send error:%v\n", err)
	}

	room.Clients[ws] = true
	fmt.Println(room.Clients)

	for {
		err = websocket.JSON.Receive(ws, &msg)
		if err != nil {
			if err.Error() == "EOF" {
				log.Printf("EOF error:%v\n", err)
				delete(room.Clients, ws)
				exitmsg := Message{RoomID: msg.RoomID, Message: msg.Name + "が退出しました", Name: "Server"}
				broadcast <- exitmsg
				break
			}
			log.Printf("Receive error:%v\n", err)
		}

		broadcast <- msg
	}
}

func handleMessages() {
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
			err := websocket.JSON.Send(client, Message{RoomID: room.ID, Message: msg.Message, Name: msg.Name})
			if err != nil {
				log.Printf("Send error:%v\n", err)
			}
		}
	}
}
