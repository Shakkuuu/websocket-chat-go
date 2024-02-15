package main

import (
	"embed"
	"fmt"
	"io"
	"log"
	"os"
	"websocket-chat/controller"
	"websocket-chat/server"
)

// view以下の静的ファイルを変数に格納
//
//go:embed view/*
var view embed.FS

func main() {
	fmt.Println("server start.")

	// 環境変数読み込み
	port := os.Getenv("SERVERPORT")        // ポート番号
	sessionKey := os.Getenv("SESSION_KEY") // セッションキー

	// アクセスログ出力用ファイル読み込み
	f, err := os.OpenFile("log/access.log", os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("logfile os.OpenFile error:%v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	// エラーログ出力用ファイル読み込み
	errorfile, err := os.OpenFile("log/error.log", os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("logfile os.OpenFile error:%v\n", err)
		os.Exit(1)
	}
	defer errorfile.Close()

	// ログの先頭に日付時刻とファイル名、行数を表示するように設定
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// エラーログの出力先をファイルに指定
	log.SetOutput(io.MultiWriter(os.Stderr, errorfile))

	controller.SessionInit(sessionKey)
	controller.TemplateInit()
	server.Init(port, view, f) // サーバ起動
}
