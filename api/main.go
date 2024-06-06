package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
)

func Env_load() {
	// `go run .` 実行で `.env.local`を使用
	// `GO_ENV=production go run .` 実行で `.env`を使用

	if os.Getenv("GO_ENV") == "" {
		os.Setenv("GO_ENV", "local")
	} else if os.Getenv("GO_ENV") == "production" {
		os.Setenv("GO_ENV", "")
	}
	
	err := godotenv.Load(fmt.Sprintf("../.env.%s", os.Getenv("GO_ENV")))
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func init() {
	Env_load() // 環境変数読み込み
}

func main() {
	http.HandleFunc("/user/create", addUbuntuUser)
	http.HandleFunc("/user/delete", deleteUbuntuUser)
	http.HandleFunc("/mails", getUserMails)
	http.HandleFunc("/mail/create", createMailAddress)
	http.HandleFunc("/mail/delete", deleteMailAddresss)
	http.HandleFunc("/sshkey/upload", uploadSSHkey)
	http.HandleFunc("/ftp/setup", setupFTP)
	http.HandleFunc("/wordpress/install", installWordpress)
	http.HandleFunc("/wordpress/status", getWordpressStatus)

	// 用意したけどどこで使うかは未定なAPI
	http.HandleFunc("/status", getControlPanelState) // コントロールパネルの死活監視用
	
	log.Fatal(http.ListenAndServe(":3000", nil))
}
