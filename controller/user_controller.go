package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"websocket-chat/entity"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

var users = []entity.User{
	{Name: "匿名", Password: "tokumei"},
	{Name: "aaa", Password: "aaa"},
}

// POST Login処理
func Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t, err := template.ParseFiles("view/login.html")
		if err != nil {
			log.Printf("controller:26, template.ParseFiles error:%v\n", err)
			http.Error(w, "ページの読み込みに失敗しました。", http.StatusInternalServerError)
			return
		}

		err = t.Execute(w, nil)
		if err != nil {
			log.Printf("controller:39, Excute error:%v\n", err)
			http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		// POSTされたものをFormから受け取り
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		tindex, err := template.ParseFiles("view/roomtop.html")
		if err != nil {
			log.Printf("controller:26, template.ParseFiles error:%v\n", err)
			http.Error(w, "ページの読み込みに失敗しました。", http.StatusInternalServerError)
			return
		}

		tlogin, err := template.ParseFiles("view/login.html")
		if err != nil {
			log.Printf("controller:26, template.ParseFiles error:%v\n", err)
			http.Error(w, "ページの読み込みに失敗しました。", http.StatusInternalServerError)
			return
		}

		for _, v := range users {
			if v.Name == username {
				if v.Password == password {
					// Room一覧をテンプレートに渡す
					var data entity.Data
					for k := range rooms {
						data.Rooms = append(data.Rooms, k)
					}

					data.Message = "ログインに成功しました。"

					session, _ := store.Get(r, "Shakku")
					session.Values["username"] = username
					session.Save(r, w)

					err = tindex.Execute(w, data)
					if err != nil {
						log.Printf("controller:39, Excute error:%v\n", err)
						http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
						return
					}
					return
				} else {
					// Room一覧をテンプレートに渡す
					var data entity.Data
					data.Message = "パスワードが違います"

					err = tlogin.Execute(w, data)
					if err != nil {
						log.Printf("controller:39, Excute error:%v\n", err)
						http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
						return
					}
					return
				}
			}
		}

		// Room一覧をテンプレートに渡す
		var data entity.Data
		data.Message = "ユーザーが存在しませんでした。"

		err = tlogin.Execute(w, data)
		if err != nil {
			log.Printf("controller:39, Excute error:%v\n", err)
			http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
			return
		}
	default:
		fmt.Fprintln(w, "controller:98, Method not allowed")
		http.Error(w, "そのメソッドは許可されていません。", http.StatusMethodNotAllowed)
		return
	}
}

// 参加ユーザーの一覧を返す
func RoomUsersList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// クエリ読み取り
		r.ParseForm()
		roomid := r.URL.Query().Get("roomid")

		var roomuserslist entity.SentRoomUsersList

		roomuserslist.UsersList = append(roomuserslist.UsersList, "匿名")

		// roomがあるか確認
		room, exists := rooms[roomid]
		if !exists {
			log.Println("controller:269, Roomが見つかりませんでした")
			http.Error(w, "Roomが見つかりませんでした", http.StatusNotFound)
			return
		}

		// Room内のユーザーを格納
		for _, user := range room.Clients {
			roomuserslist.UsersList = append(roomuserslist.UsersList, user)
		}

		// jsonに変換
		sentjson, err := json.Marshal(roomuserslist)
		if err != nil {
			log.Printf("controller:282, json.Marshal error: %v", err)
			http.Error(w, "json.Marshal error", http.StatusInternalServerError)
			return
		}

		// jsonで送信
		w.Header().Set("Content-Type", "application/json")
		w.Write(sentjson)

	default:
		fmt.Fprintln(w, "controller:292, Method not allowed")
		http.Error(w, "そのメソッドは許可されていません。", http.StatusMethodNotAllowed)
		return
	}
}

// 自身のユーザー名を返す
func GetUserName(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var sentuser entity.SentUser

		// セッション読み取り
		session, err := store.Get(r, "Shakku")
		if err != nil {
			log.Printf("controller:164, store.Get error: %v", err)
			http.Error(w, "store.Get error", http.StatusInternalServerError)
			return
		}

		username := session.Values["username"]
		if username == nil {
			fmt.Println("セッションなし")
			sentuser.Name = ""
			// jsonに変換
			sentjson, err := json.Marshal(sentuser)
			if err != nil {
				log.Printf("controller:282, json.Marshal error: %v", err)
				http.Error(w, "json.Marshal error", http.StatusInternalServerError)
				return
			}

			// jsonで送信
			w.Header().Set("Content-Type", "application/json")
			w.Write(sentjson)
			return
		}

		sentuser.Name = username.(string)

		// jsonに変換
		sentjson, err := json.Marshal(sentuser)
		if err != nil {
			log.Printf("controller:282, json.Marshal error: %v", err)
			http.Error(w, "json.Marshal error", http.StatusInternalServerError)
			return
		}

		// jsonで送信
		w.Header().Set("Content-Type", "application/json")
		w.Write(sentjson)
	default:
		fmt.Fprintln(w, "controller:292, Method not allowed")
		http.Error(w, "そのメソッドは許可されていません。", http.StatusMethodNotAllowed)
		return
	}
}
