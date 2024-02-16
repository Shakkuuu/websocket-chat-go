package main

import (
	"embed"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"websocket-chat/controller"
	"websocket-chat/db"
	"websocket-chat/server"
)

// view以下の静的ファイルを変数に格納
//
//go:embed view/*
var view embed.FS

func main() {
	// 環境変数読み込み
	connect, port, sessionKey := loadEnv()

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

	db.Init(connect)
	controller.SessionInit(sessionKey)
	controller.TemplateInit()

	fmt.Println("server start")

	go server.Init(port, view, f) // サーバ起動

	// 終了シグナルをキャッチ
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-sig
	fmt.Printf("Signal %s\n", s.String())

	// 終了処理
	fmt.Println("Shutting down server...")
	db.Close()
	fmt.Println("server stop")
}

func loadEnv() (string, string, string) {
	// Docker-compose.ymlでDocker起動時に設定した環境変数の取得
	connect := os.Getenv("DB_CONNECT")

	port := os.Getenv("SERVERPORT")        // ポート番号
	sessionKey := os.Getenv("SESSION_KEY") // セッションキー

	return connect, port, sessionKey
}
