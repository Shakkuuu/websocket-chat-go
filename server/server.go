package server

import (
	"embed"
	"log"
	"net/http"
	"os"
	"time"

	"websocket-chat/controller"

	"golang.org/x/net/websocket"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf("・[%s] %s %s %s\n", r.Method, r.RemoteAddr, r.URL, time.Since(start))
	})
}

func Init(port string, view embed.FS) {
	http.Handle("/view/", http.FileServer(http.FS(view)))

	http.Handle("/", loggingMiddleware(http.HandlerFunc(controller.Index)))              // indexページ
	http.Handle("/room", loggingMiddleware(http.HandlerFunc(controller.Room)))           // Room内のページ
	http.Handle("/rooms", loggingMiddleware(http.HandlerFunc(controller.RoomsList)))     // Room一覧取得
	http.Handle("/users", loggingMiddleware(http.HandlerFunc(controller.RoomUsersList))) // User一覧取得
	http.Handle("/ws", websocket.Handler(controller.HandleConnection))                   // メッセージWebsocket用
	go controller.HandleMessages()                                                       // goroutineとチャネルで常にメッセージを待つ

	// サーバ起動
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Printf("server:27, ListenAndServe error:%v\n", err)
		os.Exit(1)
	}
}
