package main

import (
	"embed"
	"log"
	"os"
	"websocket-chat/controller"
	"websocket-chat/server"

	"github.com/joho/godotenv"
)

// view以下の静的ファイルを変数に格納
//
//go:embed view/*
var view embed.FS

var err error

func main() {
	// ログの先頭に日付時刻とファイル名、行数を表示するように設定
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 環境変数読み込み
	err = godotenv.Load("shakku-websocket-chat.env")
	if err != nil {
		log.Printf("godotenv.Load error:%v\n", err)
		os.Exit(1)
	}
	port := os.Getenv("SERVERPORT")        // ポート番号
	sessionKey := os.Getenv("SESSION_KEY") // セッションキー
	controller.SessionInit(sessionKey)
	controller.TemplateInit()
	server.Init(port, view) // サーバ起動
}
