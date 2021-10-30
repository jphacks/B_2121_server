# goyotashi - サーバーサイド (B_2121)

## READMEはクライアントサイドのリポジトリにあります
こちらのリポジトリでは**サーバーサイドの開発技術のみ**記載していますので、製品説明等はクライアントサイドのREADMEをご覧ください。
クライアントサイドは[こちら](https://github.com/jphacks/B_2121_client)


## 開発技術

クライアントサイドは[こちら](https://github.com/jphacks/B_2121_client)

### 活用した技術

* Go
* docker
* docker-compose
* MySQL
* GitHub Actions
* Open API

#### API・データ
* ホットペッパーAPI

#### フレームワーク・ライブラリ・モジュール
* golang-migrate
* echo
* SQL Boiler

### 独自技術
#### ハッカソンで開発した独自機能・技術

* アプリ、API、BDやAPIの設計、全てこの期間で開発しています

<!--
#### 製品に取り入れた研究内容（データ・ソフトウェアなど）（※アカデミック部門の場合のみ提出必須）
* 
* 
-->

## ドキュメントなど
[server side Wiki](https://github.com/jphacks/B_2121_server/wiki)  
[client side Wiki](https://github.com/jphacks/B_2121_client/wiki)

[発表資料](https://docs.google.com/presentation/d/1oU93MItpDkqEni_x4t5PMh3QPij3ZHZhbPva-IOwYPQ/)

- サーバーサイドの docker image は https://hub.docker.com/r/kmconner/goyotashi/ から参照できる

## DB
- マイグレーションはgolang-migrateを使用している。
  - `migrate create -ext sql -dir migrations [migration name]`で新しいマイグレーションファイルを追加
  - サーバーの起動時に自動でマイグレーションされる
- `models_gen`はSQLBoilerで生成している
  - `docker compose up`でサーバーとMySQLをローカルで起動する
  - `sqlboiler mysql`で`models_gen`を生成する

## 実行方法 (docker を使う場合)

1. ホットペッパーの API キーを取得する (See https://webservice.recruit.co.jp/doc/hotpepper/reference.html)
2. `docker-compose.override.yaml.example` をコピーして `docker-compose.override.yaml` というファイル名のファイルを作成し、そこに記入する
3. `docker-compose up` (必要に応じて `--build` をつける)
