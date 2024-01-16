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
	Rooms []string
}

var rooms = make(map[string]*ChatRoom)

// var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

type Message struct {
	RoomID  string `json:"roomID"`
	Message string `json:"message"`
	Name    string `json:"name"`
}

func main() {
	http.HandleFunc("/", index)
	http.Handle("/ws", websocket.Handler(handleConnection))
	go handleMessages()

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Printf("ListenAndServe error:%v\n", err)
		os.Exit(1)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		log.Printf("template.ParseFiles error:%v\n", err)
	}

	var data Data
	for k := range rooms {
		data.Rooms = append(data.Rooms, k)
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Printf("Excute error:%v\n", err)
	}
}

func createRoom() *ChatRoom {
	roomID := fmt.Sprintf("%d", len(rooms)+1)
	room := &ChatRoom{
		ID:      roomID,
		Clients: make(map[*websocket.Conn]bool),
	}
	fmt.Printf("room %v が作成されました\n", room.ID)
	rooms[roomID] = room
	return room
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
		room = createRoom()
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
				// clients[ws] = false
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
