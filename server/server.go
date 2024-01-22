package server

import (
	"embed"
	"log"
	"net/http"
	"os"

	"websocket-chat/controller"

	"golang.org/x/net/websocket"
)

func Init(port string, view embed.FS) {
	http.Handle("/view/", http.FileServer(http.FS(view)))

	http.HandleFunc("/", controller.Index)                             // indexページ
	http.HandleFunc("/room", controller.Room)                          // Room内のページ
	http.HandleFunc("/rooms", controller.RoomsList)                    // Room一覧取得
	http.HandleFunc("/users", controller.RoomUsersList)                // User一覧取得
	http.Handle("/ws", websocket.Handler(controller.HandleConnection)) // メッセージWebsocket用
	go controller.HandleMessages()                                     // goroutineとチャネルで常にメッセージを待つ

	// サーバ起動
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Printf("ListenAndServe error:%v\n", err)
		os.Exit(1)
	}
}
