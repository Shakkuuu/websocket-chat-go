package server

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"websocket-chat/controller"

	"golang.org/x/net/websocket"
)

var err error
var accesslogfile *os.File

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		accesslog := fmt.Sprintf("%s: [%s] %s %s %s\n", timeToStr(start), r.Method, r.RemoteAddr, r.URL, time.Since(start))
		fmt.Print(accesslog)
		fmt.Fprint(accesslogfile, accesslog)
	})
}

// "YYYY-MM-DD HH-MM-SS"に変換
func timeToStr(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func Init(port string, view embed.FS, f *os.File) {
	// 出力先ファイル設定
	accesslogfile = f

	http.Handle("/view/", http.FileServer(http.FS(view)))

	http.Handle("/", loggingMiddleware(http.HandlerFunc(controller.RoomTop)))                          // roomtopページ
	http.Handle("/usermenu", loggingMiddleware(http.HandlerFunc(controller.UserMenu)))                 // usermenuページ
	http.Handle("/login", loggingMiddleware(http.HandlerFunc(controller.Login)))                       // ログインページ
	http.Handle("/signup", loggingMiddleware(http.HandlerFunc(controller.Signup)))                     // サインアップページ
	http.Handle("/logout", loggingMiddleware(http.HandlerFunc(controller.Logout)))                     // ログアウト処理
	http.Handle("/room", loggingMiddleware(http.HandlerFunc(controller.Room)))                         // Room内のページ
	http.Handle("/deleteroom", loggingMiddleware(http.HandlerFunc(controller.DeleteRoom)))             // Room削除
	http.Handle("/deleteuser", loggingMiddleware(http.HandlerFunc(controller.DeleteUser)))             // User削除
	http.Handle("/changepassword", loggingMiddleware(http.HandlerFunc(controller.ChangeUserPassword))) // パスワード更新
	http.Handle("/rooms", loggingMiddleware(http.HandlerFunc(controller.RoomsList)))                   // Room一覧取得
	http.Handle("/joinrooms", loggingMiddleware(http.HandlerFunc(controller.JoinRoomsList)))           // 参加中のRoom一覧取得
	http.Handle("/username", loggingMiddleware(http.HandlerFunc(controller.GetUserName)))              // 自身のユーザー名取得
	http.Handle("/ws", websocket.Handler(controller.HandleConnection))                                 // メッセージWebsocket用
	go controller.HandleMessages()                                                                     // goroutineとチャネルで常にメッセージを待つ

	// サーバ起動
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Printf("ListenAndServe error:%v\n", err)
		os.Exit(1)
	}
}
