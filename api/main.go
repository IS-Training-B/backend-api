package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// `json:"id"`の記載がなくてもJSON形式で返却されるが、キー名のままだと困る場合はスネークケースなどに直すことができる
type Page struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Url   string `json:"url"`
}

// Pageの配列リテラル（本来はDBから返却された値で埋めていく）
var pages = []Page {
	{
		ID:    1,
		Title: "The Go Programming Language",
		Url:   "https://golang.org/",
	}, {
		ID:    2,
		Title: "A Tour of Go",
		Url:   "https://go-tour-jp.appspot.com/welcome/1",
	},
}

// JSON返却用の構造体
type PageJSON struct {
	Status int `json:"status"`
	Pages  *[]Page
}

// ただ文字列「Hello, World」を返却するハンドラー

// func indexHandler(w http.ResponseWriter, r *http.Request) {
// 	_, err := fmt.Fprint(w, "Hello, World")
// 	if err != nil {
// 		return
// 	}
// }

// appilication/jsonでJSONっぽい値を返却するハンドラー
func pagesHandler(w http.ResponseWriter, r *http.Request) {

	var pj PageJSON
	pj.Status = 200
	pj.Pages = &pages

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(&pj); err != nil {
		log.Fatal(err)
	}
	fmt.Println(buf.String())

	// Content-Typeを設定
	w.Header().Set("Content-Type", "application/json;charset=utf-8")

	// Responseに書き込み
	_, err := fmt.Fprint(w, buf.String())
	if err != nil {
		return
	}
}

func main() {

	// text/plain返却　"/" のとき indexHandlerを実行する
	// http.HandleFunc("/", indexHandler)

	// application/json返却 "/pages"の時 pagesHandlerを返却
	http.HandleFunc("/pages", pagesHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
	
}