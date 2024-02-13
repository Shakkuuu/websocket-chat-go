package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"websocket-chat/entity"
)

// var store = sessions.NewCookieStore([]byte("something-very-secret"))

// TODO:mapにしてみる？
var users = []entity.User{
	{Name: "匿名", Password: "qawsedrftgyhujikolp"},
}

// Signup処理
func Signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err = tsignup.Execute(w, nil)
		if err != nil {
			log.Printf("Excute error:%v\n", err)
			http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		// POSTされたものをFormから受け取り
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")
		checkpass := r.FormValue("checkpassword")

		if username == "" || password == "" || checkpass == "" {
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "入力されていない項目があります。"

			err = tsignup.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		if password != checkpass {
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "確認用再入力パスワードが一致していません。"

			err = tsignup.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		for _, v := range users {
			if v.Name == username {
				// メッセージをテンプレートに渡す
				var data entity.Data
				data.Message = "そのユーザー名は既に使用されています。"

				err = tsignup.Execute(w, data)
				if err != nil {
					log.Printf("Excute error:%v\n", err)
					http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
					return
				}
				return
			}
		}

		user := entity.User{
			Name:               username,
			Password:           password,
			ParticipatingRooms: make(map[*entity.ChatRoom]bool),
		}

		users = append(users, user)

		// メッセージをテンプレートに渡す
		var data entity.Data
		data.Message = "登録が完了しました。ログインしてください。"

		err = tlogin.Execute(w, data)
		if err != nil {
			log.Printf("Excute error:%v\n", err)
			http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
			return
		}
	default:
		fmt.Fprintln(w, "Method not allowed")
		http.Error(w, "そのメソッドは許可されていません。", http.StatusMethodNotAllowed)
		return
	}
}

// Login処理
func Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err = tlogin.Execute(w, nil)
		if err != nil {
			log.Printf("Excute error:%v\n", err)
			http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		// POSTされたものをFormから受け取り
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "入力されていない項目があります。"

			err = tlogin.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		for _, v := range users {
			if v.Name == username {
				if v.Password == password {
					// Room一覧とメッセージをテンプレートに渡す
					var data entity.Data
					for k := range rooms {
						data.Rooms = append(data.Rooms, k)
					}

					data.Message = "ログインに成功しました。"

					session, _ = store.Get(r, SESSION_NAME)
					session.Values["username"] = username
					session.Save(r, w)

					err = troomtop.Execute(w, data)
					if err != nil {
						log.Printf("Excute error:%v\n", err)
						http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
						return
					}
					return
				} else {
					// メッセージをテンプレートに渡す
					var data entity.Data
					data.Message = "パスワードが違います"

					err = tlogin.Execute(w, data)
					if err != nil {
						log.Printf("Excute error:%v\n", err)
						http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
						return
					}
					return
				}
			}
		}

		// メッセージをテンプレートに渡す
		var data entity.Data
		data.Message = "ユーザーが存在しませんでした。"

		err = tlogin.Execute(w, data)
		if err != nil {
			log.Printf("Excute error:%v\n", err)
			http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
			return
		}
	default:
		fmt.Fprintln(w, "Method not allowed")
		http.Error(w, "そのメソッドは許可されていません。", http.StatusMethodNotAllowed)
		return
	}
}

// Logout処理
func Logout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		session, err = store.Get(r, SESSION_NAME)
		if err != nil {
			log.Printf("store.Get error: %v", err)
			http.Error(w, "store.Get error", http.StatusInternalServerError)
			return
		}
		// セッション削除
		session.Options.MaxAge = -1
		session.Save(r, w)

		// メッセージをテンプレートに渡す
		var data entity.Data
		data.Message = "ログアウトしました。ログインしてください。"

		err = tlogin.Execute(w, data)
		if err != nil {
			log.Printf("Excute error:%v\n", err)
			http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
			return
		}
	default:
		fmt.Fprintln(w, "Method not allowed")
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
			log.Println("Roomが見つかりませんでした")
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
			log.Printf("json.Marshal error: %v", err)
			http.Error(w, "json.Marshal error", http.StatusInternalServerError)
			return
		}

		// jsonで送信
		w.Header().Set("Content-Type", "application/json")
		w.Write(sentjson)

	default:
		fmt.Fprintln(w, "Method not allowed")
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
		session, err := store.Get(r, SESSION_NAME)
		if err != nil {
			log.Printf("store.Get error: %v", err)
			http.Error(w, "store.Get error", http.StatusInternalServerError)
			return
		}

		username := session.Values["username"]
		if username == nil {
			fmt.Println("セッションなし")
			http.Error(w, "session not found", http.StatusUnauthorized)
			return
		}

		sentuser.Name = username.(string)

		// jsonに変換
		sentjson, err := json.Marshal(sentuser)
		if err != nil {
			log.Printf("json.Marshal error: %v", err)
			http.Error(w, "json.Marshal error", http.StatusInternalServerError)
			return
		}

		// jsonで送信
		w.Header().Set("Content-Type", "application/json")
		w.Write(sentjson)
	default:
		fmt.Fprintln(w, "Method not allowed")
		http.Error(w, "そのメソッドは許可されていません。", http.StatusMethodNotAllowed)
		return
	}
}
