package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"websocket-chat/entity"
	"websocket-chat/model"

	"gorm.io/gorm"
)

// ユーザー情報のページ
func UserMenu(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// セッション読み取り
		un, err := SessionToGetName(r)
		if err != nil {
			log.Printf("SessionToGetName error: %v", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "再ログインしてください"

			err = tlogin.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		user, err := model.GetUserByName(un)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("model.GetUserByName error: %v", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "ユーザーが見つかりませんでした。"

			err = tlogin.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}
		if err != nil {
			log.Printf("GetUserByName error: %v", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "データベースとの接続に失敗しました。"

			err = tlogin.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// メッセージをテンプレートに渡す
		var data entity.Data
		data.Message = user.Name + "さん、こんにちは。"

		data.Name = user.Name

		err = tusermenu.Execute(w, data)
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

// ユーザー削除
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// セッション読み取り
		session, err = store.Get(r, SESSION_NAME)
		if err != nil {
			log.Printf("store.Get error: %v", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "再ログインしてください"

			err = tlogin.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		username := session.Values["username"]
		if err != nil {
			log.Printf("Session.Values error: %v", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "再ログインしてください"

			err = tlogin.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}
		un := username.(string)

		var user entity.User
		// セッションのユーザー取得
		user, err = model.GetUserByName(un)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("model.GetUserByName error: %v", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "ユーザーが見つかりませんでした。"

			err = tlogin.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}
		if err != nil {
			log.Printf("model.GetUserByName error: %v", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "データベースとの接続に失敗しました。"

			err = tusermenu.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// セッション削除
		session.Options.MaxAge = -1
		session.Save(r, w)

		// ユーザーが作成したRoomの削除
		var proom entity.ParticipatingRoom
		prooms, err := model.GetParticipatingRoom(user.Name)
		if err != nil {
			log.Printf("model.GetParticipatingRoom error: %v", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "データベースとの接続に失敗しました。"

			err = tusermenu.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		for _, proom = range prooms {
			if !proom.IsMaster {
				continue
			}

			// ユーザーの参加中ルームリストからも削除
			err = model.DeleteParticipatingRoomByRoomID(proom.RoomID)
			if err != nil {
				log.Printf("model.DeleteParticipatingRoomByRoomID error: %v", err)
				// メッセージをテンプレートに渡す
				var data entity.Data
				data.Message = "データベースとの接続に失敗しました。"

				err = tusermenu.Execute(w, data)
				if err != nil {
					log.Printf("Excute error:%v\n", err)
					http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
					return
				}
				return
			}

			// 部屋削除
			err = model.DeleteRoom(proom.RoomID)
			if err != nil {
				log.Printf("model.DeleteRoom error: %v", err)
				// メッセージをテンプレートに渡す
				var data entity.Data
				data.Message = "データベースとの接続に失敗しました。"

				err = tusermenu.Execute(w, data)
				if err != nil {
					log.Printf("Excute error:%v\n", err)
					http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
					return
				}
				return
			}
		}

		// ユーザーの参加中ルームリストを削除
		err = model.DeleteParticipatingRoomByUserName(user.Name)
		if err != nil {
			log.Printf("model.DeleteParticipatingRoomUserName error: %v", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "データベースとの接続に失敗しました。"

			err = tusermenu.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// ユーザー削除
		err = model.DeleteUser(user.Name)
		if err != nil {
			log.Printf("model.DeleteUser error: %v", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "データベースとの接続に失敗しました。"

			err = tusermenu.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// メッセージをテンプレートに渡す
		var data entity.Data
		data.Message = "ユーザーを削除しました。"

		err = tsignup.Execute(w, data)
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

// パスワード変更
func ChangeUserPassword(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		r.ParseForm()
		oldpassword := r.FormValue("oldpassword")
		password := r.FormValue("password")
		checkpass := r.FormValue("checkpassword")

		// セッション読み取り
		session, err = store.Get(r, SESSION_NAME)
		if err != nil {
			log.Printf("store.Get error: %v", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "再ログインしてください"

			err = tlogin.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		username := session.Values["username"]
		if err != nil {
			log.Printf("Session.Values error: %v", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "再ログインしてください"

			err = tlogin.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}
		un := username.(string)

		var user entity.User
		// セッションのユーザー取得
		user, err = model.GetUserByName(un)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("model.GetUserByName error: %v", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "ユーザーが見つかりませんでした。"

			err = tlogin.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}
		if err != nil {
			log.Printf("model.GetUserByName error: %v", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "データベースとの接続に失敗しました。"

			err = tusermenu.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		if oldpassword == "" || password == "" || checkpass == "" {
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "入力されていない項目があります。"

			err = tusermenu.Execute(w, data)
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

			err = tusermenu.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		err = model.HashPassCheck(user.Password, oldpassword)
		if err != nil {
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "現在のパスワードが違います"

			err = tusermenu.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}
		hashpass, err := model.HashPass(password)
		if err != nil {
			log.Printf("bcrypt.GenerateFromPassword error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "パスワードのハッシュに失敗しました。"

			err = tusermenu.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		user = entity.User{
			ID:       user.ID,
			Name:     user.Name,
			Password: hashpass,
		}

		// ユーザー更新
		err = model.PutUserByName(&user, un)
		if err != nil {
			log.Printf("model.PutUserByName error: %v", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "データベースとの接続に失敗しました。"

			err = tusermenu.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// 再ログイン用に一度セッション削除
		session.Options.MaxAge = -1
		session.Save(r, w)

		// メッセージをテンプレートに渡す
		var data entity.Data
		data.Message = "パスワードを更新しました。再ログインしてください。"

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

		// ユーザー一覧取得
		users, err := model.GetUsers()
		if err != nil {
			log.Printf("model.GetUsers error: %v", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "データベースとの接続に失敗しました。"

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

		hashpass, err := model.HashPass(password)
		if err != nil {
			log.Printf("bcrypt.GenerateFromPassword error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "パスワードのハッシュに失敗しました。"

			err = tsignup.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		user := entity.User{
			Name:     username,
			Password: hashpass,
		}

		// ユーザー追加
		err = model.AddUser(&user)
		if err != nil {
			log.Printf("model.AddUser error: %v", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "データベースとの接続に失敗しました。"

			err = tsignup.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

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

		var user entity.User
		// 登録されているユーザー取得
		user, err = model.GetUserByName(username)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("model.GetUserByName error: %v", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "ユーザーが見つかりませんでした。"

			err = tlogin.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}
		if err != nil {
			log.Printf("model.GetUserByName error: %v", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "データベースとの接続に失敗しました。"

			err = tlogin.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// Room一覧取得
		rooms := model.GetRooms()

		err = model.HashPassCheck(user.Password, password)
		if err == nil {
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
		// セッション確認
		session, err = store.Get(r, SESSION_NAME)
		if err != nil {
			log.Printf("store.Get error: %v", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "再ログインしてください"

			err = tlogin.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
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

// 自身のユーザー名を返す
func GetUserName(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var sentuser entity.SentUser

		// セッション読み取り
		un, err := SessionToGetName(r)
		if err != nil {
			log.Printf("SessionToGetName error: %v", err)
			log.Println("セッションが見つかりませんでした")
			http.Error(w, "セッションが見つかりませんでした", http.StatusNotFound)
			return
		}

		sentuser.Name = un

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

// 参加中のRoomの一覧を返す
func JoinRoomsList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var joinroomslist entity.SentRoomsList

		// セッション読み取り
		un, err := SessionToGetName(r)
		if err != nil {
			log.Printf("SessionToGetName error: %v", err)
			// メッセージをテンプレートに渡す
			var data entity.Data
			data.Message = "再ログインしてください"

			err = tlogin.Execute(w, data)
			if err != nil {
				log.Printf("Excute error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		var user entity.User
		// ユーザー取得
		user, err = model.GetUserByName(un)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("model.GetUserByName error: %v", err)
			log.Printf("User Not Found: %v", err)
			return
		}
		if err != nil {
			log.Printf("GetUserByName error: %v", err)
			return
		}

		// joinRoomを格納
		prooms, err := model.GetParticipatingRoom(user.Name)
		if err != nil {
			fmt.Println("データベースとの接続に失敗しました。")
			http.Error(w, "GetParticipatingRoom error", http.StatusUnauthorized)
			return
		}
		for _, proom := range prooms {
			joinroomslist.RoomsList = append(joinroomslist.RoomsList, proom.RoomID)
		}

		fmt.Println(joinroomslist.RoomsList)

		// jsonに変換
		sentjson, err := json.Marshal(joinroomslist)
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
