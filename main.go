package main

import (
	"log"
	"os"
	"websocket-chat/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("godotenv.Load error:%v\n", err)
		os.Exit(1)
	}
	port := os.Getenv("SERVERPORT")
	server.Init(port)
}
