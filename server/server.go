package server

import (
	"log"
	"net/http"
	"os"

	"websocket-chat/controller"

	"golang.org/x/net/websocket"
)

func Init(port string) {
	http.HandleFunc("/", controller.Index)                             // indexページ
	http.HandleFunc("/room", controller.Room)                          // Room内のページ
	http.Handle("/ws", websocket.Handler(controller.HandleConnection)) // メッセージWebsocket用
	go controller.HandleMessages()                                     // goroutineとチャネルで常にメッセージを待つ

	// サーバ起動
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Printf("ListenAndServe error:%v\n", err)
		os.Exit(1)
	}
}
