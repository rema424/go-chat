# go-chat

## memo

### 1. プロジェクトのの初期化

#### 1.1. リポジトリの作成

#### 1.2. go.mod の作成

```sh
go mod init chat
```

### 2. auth パッケージ

#### 2.1. データベースの作成

```sh
mysql.server start
mysql -uroot -e 'create database if not exists go_chat;'
mysql -uroot -e 'create user if not exists devuser@localhost identified by "Passw0rd!";'
mysql -uroot -e 'grant all privileges on go_chat.* to devuser@localhost;'
mysql -uroot -e 'show databases;'
mysql -uroot -e 'select host, user from mysql.user;'
mysql -uroot -e 'show grants for devuser@localhost;'
```

#### 2.2 テスト

```sh
go test ./... -cover -coverprofile=cover.out
go tool cover -html=cover.out
```
