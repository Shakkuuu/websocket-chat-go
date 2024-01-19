package main

import (
	"log"
	"os"
	"websocket-chat/server"

	"github.com/joho/godotenv"
)

func main() {
	// 環境変数読み込み
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("godotenv.Load error:%v\n", err)
		os.Exit(1)
	}
	port := os.Getenv("SERVERPORT") // ポート番号
	server.Init(port)               // サーバ起動
}
