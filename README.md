# websocket-chat-go

## 注意

- 以下の内容を含んだ"〇〇.env"ファイルを配置してください。

```:env
SERVERPORT=":8000"
SESSION_KEY="hogefuga"
```

- logディレクトリ内に`access.log`と`error.log`というファイルを配置してください。

## やること

- 見た目
- コードまとめる
- ユーザー一覧をdbに保存(mysqlじゃない簡易的なやつでいいかも)
- ユーザー削除
- パスワード変更
- ユーザー名変更
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
- 人がいなくなったら自動削除？ no
- クライアントがわのName消せるかもそれでサーバー側のSessionのみでName管理 wsの時につまる
- ログイン画面にゲストログイン追加 nasi
- ユーザー作成時にユーザーID割り振り nasi
- ユーザーIDをルームに追加していく nasi

- 参加ユーザー名を押して、プライベートチャット送信先を選択できるようにする？ やっぱやる？
