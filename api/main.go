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
	http.HandleFunc("/mails", getUserMails)
	http.HandleFunc("/mail/create", createMailAddress)
	http.HandleFunc("/mail/delete", deleteMailAddresss)
	http.HandleFunc("/status", getControlPanelState)
	log.Fatal(http.ListenAndServe(":3000", nil))
}