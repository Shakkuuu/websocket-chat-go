package server

import (
	"log"
	"net/http"
	"os"

	"websocket-chat/controller"

	"golang.org/x/net/websocket"
)

func Init() {
	http.HandleFunc("/", controller.Index)
	http.HandleFunc("/room", controller.Room)
	http.Handle("/ws", websocket.Handler(controller.HandleConnection))
	go controller.HandleMessages()

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Printf("ListenAndServe error:%v\n", err)
		os.Exit(1)
	}
}
