# websocket-chat-go

## 概要

WebosocketとGo言語を主に使用した、リアルタイム性のあるチャットアプリケーション。

ユーザー登録し、Roomを作成したり参加したりすることで、自由にさまざまなユーザーとリアルタイムに会話することができる。

ユーザーを指定することで、そのRoom内で指定したユーザーにのみ表示されるプライベートメッセージを送信することもできる。

サービスは[Render](https://render.com/)というPassSサービスにデプロイしており、ユーザーやRoomのDB保存はRender内のPostgreSQLサービスを使用している。

デプロイしたサービスのURLは[こちら](https://shakku-websocket-chat.onrender.com/)。

プログラムの細かい解説はQiitaの[この記事](https://qiita.com/Shakku)をご覧ください。

※2024/02/20 記事未作成。

## ディレクトリ構成

```c
websocket-chat-go/
├── src
│   ├── controller
│   │   ├── room_controller.go
│   │   ├── session.go
│   │   ├── template.go
│   │   └── user_controller.go
│   ├── db
│   │   └── db.go
│   ├── entity
│   │   └── entity.go
│   ├── log
│   │   ├── access.log
│   │   ├── error.log
│   │   └── chat.log
│   ├── model
│   │   ├── room_model.go
│   │   └── user_model.go
│   ├── server
│   │   └── server.go
│   ├── view
│   │   ├── icon
│   │   │   └── favicon.ico
│   │   ├── script
│   │   │   ├──room.js
│   │   │   ├──roomtop.js
│   │   │   ├──signup-login.js
│   │   │   └── usermenu.js
│   │   ├── style
│   │   │   └── main.css
│   │   ├── login.html
│   │   ├── room.html
│   │   ├── roomtop.html
│   │   ├── signup.html
│   │   └── usermenu.html
│   ├── dockerfile
│   ├── go.mod
│   ├── go.sum
│   └── main.go
├── .env(開発用)
├── .gitignore
├── docker-compose.yml
├── Production.env(デプロイ用)
└── README.md
```

## 起動

### 開発時

docker-composeを使用して開発。

.envファイルに以下の環境変数を記載しておく。

```env
SERVERPORT=Goのサーバーのポート(今回は":8000"とした)
SESSION_KEY=セッションキー
DB_PROTOCOL=DB接続用のホスト名(docker-composeで起動する場合は、サービス名の"db"とする。)
DB_DATABASENAME=データベース名
DB_USERNAME=データベース用のユーザー名
DB_USERPASS=データベース用のパスワード
DB_PORT=データベースのポート(今回は"5432")


SERVERPORT=":8000"
SESSION_KEY="shakku-session-key"
DB_PROTOCOL="dpg-cn7ac5821fec73fkpt9g-a"
DB_DATABASENAME="shakku_websocket_chat_sql"
DB_USERNAME="shakku_websocket_chat_sql_user"
DB_USERPASS="N9OH1wGu4f5Khzv9xPoEY8SoCbKs0mhx"
DB_PORT="5432"
PORT="8000"
```

起動

```docker
docker compose up -d
```

ポート

- app:8000
- db:5432

### デプロイ時

[Render](https://render.com/)というPassSサービスを使用した。

デプロイ設定時に環境変数を入力する。

```env
SERVERPORT=Goのサーバーのポート(今回は":8000"とした)
SESSION_KEY=セッションキー
DB_PROTOCOL=DB接続用のホスト名(docker-composeで起動する場合は、サービス名の"db"とする。)
DB_DATABASENAME=データベース名
DB_USERNAME=データベース用のユーザー名
DB_USERPASS=データベース用のパスワード
DB_PORT=データベースのポート(今回は"5432")
PORT=RenderでのGoサーバー起動用に指定するポート(今回は"8000")
```

詳細な起動方法はQiita記事に記載。

ポート

- app:8000
- db:5432

## 使用した技術

- Golang
- PostgreSQL
- GORM
- Websocket
- Session
- Cookie
- bcrypt
- Render
- Javascript
- CSS
- Docker
- Docker Compose

### 使用したパッケージ

- embed
- errors
- fmt
- io
- log
- os
- os/signal
- syscall
- time
- encoding/json
- net/http
- html/template
- strings
- golang.org/x/net/websocket v0.20.0
- github.com/gorilla/sessions v1.2.2
- golang.org/x/crypto v0.19.0
- gorm.io/driver/postgres v1.5.6
- gorm.io/gorm v1.25.7

## ルーティング

### ユーザー系

- SignUpページ

`GET /signup`

- Signup

`POST /signup`

- Loginページ

`GET /login`

- Login

`POST /login`

- Logout

`GET /logout`

- ユーザーのメニューページ

`GET /usermenu`

- ユーザー削除

`GET /deleteuser`

- パスワード変更

`POST /changepassword`

- 自身のユーザー名取得

`GET /username`

- そのユーザーの参加中のRoom一覧

`GET /joinrooms`

### Room系

- Roomのtopページ

`GET /`

- Room作成

`POST /`

- Room内のページ

`GET /room`

- Room削除

`GET /deleteroom`

- Room一覧取得

`GET /rooms`

- Room内のユーザー一覧取得

`GET /roomusers`

### Websocket系

- コネクション確立と、メッセージ受信待機

`/ws`

- controller.HandleMessages()

`goroutineでチャネルにクライアントから来たメッセージが入るまで待機し、クライアントにメッセージを送信する。`

## やること

- Qiita記事書く
- 見た目
- コードまとめる
- サーバーを安全に終了Graceful shutdown (Dockerfileでgo run main.goのコマンドを実行させずに、コンテナに入って実行しないと、docker-composeを終了時にコンテナ自体が終了されてしまい、goにシグナルがうまく飛ばされない)
- ログをバッファリングして、ある程度溜まってから書き込み 何行ごととか
- htmlで入力必須かつpattern入れているが、サーバー側でも実装する費用があるのか

## 済

- ユーザーログイン機能入れる？ ok
- ユーザーごとに参加サーバー一覧リスト持っている mapで実装した ok
- サーバ再起動後、クッキーが残っているときに、users内にそのユーザーがいるか確認 ok
- ログアウト機能 ok
- log.printfの行数修正 ok
- 部屋削除 ok
- 部屋作成者が部屋削除できるようにする？ ok
- パスワードに日本語入力できないように ok半角英数字をそれぞれ1種類以上含み、8文字以上100文字以下にした
- templateをサーバー起動時に一括で読み込み ok
- sessionをサーバー起動時に設定してenvに対応できるように ok
- セッションキー デプロイ時はコンテナ起動時などに設定しておく、.envでやるのはキツそう？ できた
- session取得系の処理を関数でまとめる ok
- ログをファイルに出力 ok (log.SetOutput()でログの出力先を変更できた。io.MultiWriter(os.Stderr, errorfile)を使用することで、標準エラー出力とerror.logファイルの両方にエラーログを出力できるようにした。アクセスの方はfmt.Fprintで出力)
- ユーザー一覧をdbに保存(mysql) ok
- ユーザー削除 ok
- Room一覧をdbに保存(mysql) ok
- パスワード変更 ok
- チャットログをlogファイルに残す(サーバー名、ユーザー名、宛先、メッセージ、日時を入れる) ok
- メッセージに改行が含まれるとチャットログが途中で改行されてしまう ok

- 人がいなくなったら自動削除？ no
- クライアントがわのName消せるかもそれでサーバー側のSessionのみでName管理 wsの時につまる
- ログイン画面にゲストログイン追加 nasi
- ユーザー作成時にユーザーID割り振り nasi
- ユーザーIDをルームに追加していく nasi
- ユーザー名変更 no (ユーザー名で色々管理しているから実装大変そう)

- 参加ユーザー名を押して、プライベートチャット送信先を選択できるようにする？ やっぱやる？
