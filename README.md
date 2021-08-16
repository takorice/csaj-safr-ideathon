## 概要
CSAJ 顔認証アイデアソン2020 で実装したバッチ、APIです。

- SAFR (顔認証基盤) 上に蓄積された各ユーザの感情値を取得し、DBへ格納するバッチ
  - SAFR のAPIを利用
- DBに格納された感情値をレスポンスするAPI


## 技術要素
- backend
  - Go (Gin)
  - Docker
  - PostgreSQL
  - heroku


## 実行
### docker build
```
$ make build
```

### Web と DB 起動
```
$ make up
```

http://localhost:8000/api/reaction_summaries

(上記で、main.goを起動し、テーブルをAutoMigrateする)

### Worker 実行
```
$ make task
```
(main.goを起動し、SAFRから各ユーザの直近の感情値を取得し、DBへ格納する)

### 終了
```
$ make down
```


## APIの利用
### Heroku 上の API を呼び出し、感情値を取得する

- 全件取得  

https://csaj-ideathon.herokuapp.com/api/reaction_summaries

- 期間を絞って取得  

https://csaj-ideathon.herokuapp.com/api/reaction_summaries?reacted_at_from=2020-07-02T06:28:17&reacted_at_to=2020-07-02T07:51:00
