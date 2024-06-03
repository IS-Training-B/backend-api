## ローカル環境: Go セットアップ
### Go (v1.22.3) インストール
  - Homebrewを使う場合（macOSまたはLinux）
    1. [Homebrewインストール](https://brew.sh/ja/)
    2. `brew install go@1.22.3`
  - ライブラリのインストール  
    `go get github.com/joho/godotenv`

## ローカル環境: DBセットアップ
### 1. MariaDB(v10.6)　インストール  
 - Homebrewを使う場合  
    `brew install mariadb@10.6`
### 2. PATH設定  
    `echo 'export PATH="/opt/homebrew/opt/mariadb@10.6/bin:$PATH"' >> ~/.zshrc`
### 3. rootアカウントのパスワード設定  
`sudo mariadb-secure-installation`

Enter current password for root (enter for none)：  
と表示されてもrootの初期パスワードはnullなのでそのままEnter

残りの設定はこのページ見ながら↓  
https://qiita.com/ynack/items/4709b77d42847823cdb3

### 4. SQLファイルの実行
`mysql -u root -p`
- テーブル作成
    `source　.../backend-api/setup.sql`
- テストユーザデータ追加
    `source .../backend-api/test_users.sql`
### 5. ビルド
`go run main.go db.go`  
3000ポート