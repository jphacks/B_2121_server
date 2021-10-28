# goyotashi
クライアントサイドは[こちら](https://github.com/jphacks/B_2121_client)

[![IMAGE ALT TEXT HERE](https://jphacks.com/wp-content/uploads/2021/07/JPHACKS2021_ogp.jpg)](https://www.youtube.com/watch?v=LUPQFB4QyVo)

## 製品概要
グループでよく行く飲食店のリストを共有する *eat\*Tech*（消費者の食行動\*Tech）SNS 「**goyotashi**」

### 背景(製品開発のきっかけ、課題等）
グループや個人で飲食店に行くとき、「どこでもいい」と言いながら決まらないという面倒な問題がある。  
人は既知の飲食店から行き先を選ぶ際、無意識によく行く飲食店リストを思い出そうとし、これでもない、あれでもないと悩む。  
そんなときに見る便利なリストが欲しかった。

### 製品説明（具体的な製品の説明）
定期的にご飯を食べる人達(**グループ**)でよく行くお店リストを集約・可視化し、飲食店選びのコストを下げてくれるiOSアプリ。  
リストはグループ外にも公開され、副作用として他のグループ御用達の、信頼できる飲食店を知ることができる。

### 特長
#### 1. 知っている飲食店の中から選ぶので、選択肢が多くなくて選びやすい！

#### 2. すぐに思い出せないが、今の気分にあっている飲食店が見つかる！

#### 3. 自分が属していない近くのグループや、似たような嗜好性のグループ御用達の飲食店を見ることで、選択肢が広がる！ 

### 解決出来ること
「どこでもいい」と言いながらも、潜在的には食べたいものの気分がある。  
そんなときに検索に変わるより低コストな手段として、飲食店選びを解決する。

### 今後の展望
グループが公開されることによる副作用をより活かす(ex. グループ・嗜好性情報の拡充)

### 注力したこと（こだわり等）
* 課題の吟味、深掘り
* ユーザーインタビューやユースケースの具体化
* 既存のグルメアプリとの差別化

## 開発技術
### 活用した技術
* クライアント:swift
* サーバ:Go

#### API・データ
* 
* 

#### フレームワーク・ライブラリ・モジュール
* 
* 

#### デバイス
* 

### 独自技術
#### ハッカソンで開発した独自機能・技術
* 独自で開発したものの内容をこちらに記載してください
* 特に力を入れた部分をファイルリンク、またはcommit_idを記載してください。

<!--
#### 製品に取り入れた研究内容（データ・ソフトウェアなど）（※アカデミック部門の場合のみ提出必須）
* 
* 
-->

## ドキュメントなど
[server side Wiki](https://github.com/jphacks/B_2121_server/wiki)  
[client side Wiki](https://github.com/jphacks/B_2121_client/wiki)  

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
