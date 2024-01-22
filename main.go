package main

import (
	"embed"
	"log"
	"os"
	"websocket-chat/server"

	"github.com/joho/godotenv"
)

// view以下の静的ファイルを変数に格納
//
//go:embed view/*
var view embed.FS

func main() {
	// 環境変数読み込み
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("main:21, godotenv.Load error:%v\n", err)
		os.Exit(1)
	}
	port := os.Getenv("SERVERPORT") // ポート番号
	server.Init(port, view)         // サーバ起動
}
