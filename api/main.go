package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
)

type Page struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Url   string `json:"url"`
}

// jsonの仮データ
var pages = []Page {
	{
		ID:    1,
		Title: "The Go Programming Language",
		Url:   "https://golang.org/",
	}, {
		ID:    2,
		Title: "A Tour of Go",
		Url:   "https://go-tour-jp.appspot.com/welcome/1",
	},{
		ID:    3,
		Title: "A Tour of Go",
		Url:   "https://go-tour-jp.appspot.com/welcome/1",
	}, {
		ID:    4,
		Title: "A Tour of Go",
		Url:   "https://go-tour-jp.appspot.com/welcome/1",
	}, {
		ID:    5,
		Title: "A Tour of Go",
		Url:   "https://go-tour-jp.appspot.com/welcome/1",
	},
}

type PageJSON struct {
	Status int `json:"status"`
	Pages  *[]Page
}

func Env_load() {
	// `go run main.go db.go` 実行で `.env.local`を使用
	// `GO_ENV=production go run main.go db.go` 実行で `.env`を使用
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

func main() {
	// 環境変数読み込み
	Env_load()
	// DB接続
	Db()
	
	http.HandleFunc("/pages", pagesHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func pagesHandler(w http.ResponseWriter, r *http.Request) {
	var pj PageJSON
	pj.Status = 200
	pj.Pages = &pages

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(&pj); err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")

	// Responseに書き込み
	_, err := fmt.Fprint(w, buf.String())
	if err != nil {
		return
	}
}
