package db

import (
	"fmt"
	"log"
	"time"
	"websocket-chat/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db      *gorm.DB
	err     error
	tokumei entity.User
)

// データベースと接続
func Init(host, user, password, database string) {
	// DBMS := dbms         // データベースの種類
	// USER := username     // ユーザー名
	// PASS := userpass     // パスワード
	// PROTOCOL := protocol // 3306ポート
	// DBNAME := dbname     // データベース名

	// CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"
	CONNECT := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai", host, user, password, database)

	fmt.Println(CONNECT)
	fmt.Println("DB接続開始")
	// 接続できるまで一定回数リトライ
	count := 0
	db, err = gorm.Open(postgres.Open(CONNECT), &gorm.Config{})
	if err != nil {
		for {
			if err == nil {
				fmt.Println("")
				break
			}
			fmt.Print(".")
			time.Sleep(time.Second)
			count++
			if count > 10 { // countgaが180になるまでリトライ
				fmt.Println("")
				log.Printf("db Init error: %v\n", err)
				panic(err)
			}
			db, err = gorm.Open(postgres.Open(CONNECT), &gorm.Config{})
		}
	}
	autoMigration()

	var p entity.ParticipatingRoom

	db.Delete(&p)

	var u entity.User

	db.Where("name = ?", "匿名").Delete(&u)

	insertTokumei()
	fmt.Println("DB接続完了")
}

// serviceでデータベースとやりとりする用
func GetDB() *gorm.DB {
	return db
}

// サーバ終了時にデータベースとの接続終了
func Close() {
	// var p entity.ParticipatingRoom

	// db.Delete(&p)

	if sqlDB, err := db.DB(); err != nil {
		log.Printf("db Close error: %v\n", err)
		panic(err)
	} else {
		if err := sqlDB.Close(); err != nil {
			log.Printf("db Close error: %v\n", err)
			panic(err)
		}
	}
	// if err := db.Close(); err != nil {
	// 	log.Printf("db Close error: %v\n", err)
	// 	panic(err)
	// }
}

// entityを参照してテーブル作成　マイグレーション
func autoMigration() {
	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.ParticipatingRoom{})
}

// 匿名ユーザーを初期に追加
func insertTokumei() {
	db := GetDB()
	tokumei = entity.User{
		Name:     "匿名",
		Password: "tokumei",
	}

	err := db.Create(&tokumei).Error
	if err != nil {
		log.Printf("db.Create tokumei error: %v\n", err)
	}
	log.Println("匿名ユーザーが登録されました。")
}
