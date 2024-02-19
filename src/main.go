package main

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"websocket-chat/controller"
	"websocket-chat/db"
	"websocket-chat/model"
	"websocket-chat/server"

	"gorm.io/gorm"
)

// view以下の静的ファイルを変数に格納
//
//go:embed view/*
var view embed.FS

func main() {
	// 環境変数読み込み
	host, user, password, database, port, sessionKey := loadEnv()

	// アクセスログ出力用ファイル読み込み
	f, err := os.OpenFile("log/access.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("logfile os.OpenFile error:%v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	// エラーログ出力用ファイル読み込み
	errorfile, err := os.OpenFile("log/error.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("logfile os.OpenFile error:%v\n", err)
		os.Exit(1)
	}
	defer errorfile.Close()

	// ログの先頭に日付時刻とファイル名、行数を表示するように設定
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// エラーログの出力先をファイルに指定
	log.SetOutput(io.MultiWriter(os.Stderr, errorfile))

	db.Init(host, user, password, database)
	controller.SessionInit(sessionKey)
	err = controller.TemplateInit()
	if err != nil {
		log.Printf("TemplateInit error:%v\n", err)
		os.Exit(1)
	}
	err = model.RoomInit()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("Roomが空の初期状態です。")
	}
	if err != nil {
		log.Printf("RoomInit error:%v\n", err)
		os.Exit(1)
	}

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

func loadEnv() (string, string, string, string, string, string) {
	// Docker-compose.ymlでDocker起動時に設定した環境変数の取得
	// dbms := os.Getenv("DB_DBMS")           // データベースの種類
	username := os.Getenv("DB_USERNAME")   // データベースのユーザー名
	userpass := os.Getenv("DB_USERPASS")   // データベースのユーザーのパスワード
	protocol := os.Getenv("DB_PROTOCOL")   // データベースの使用するプロトコル
	dbname := os.Getenv("DB_DATABASENAME") // データベース名

	port := os.Getenv("SERVERPORT")        // ポート番号
	sessionKey := os.Getenv("SESSION_KEY") // セッションキー

	return protocol, username, userpass, dbname, port, sessionKey
}
