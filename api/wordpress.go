package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// POST localhost:3000/wordpress/install
func installWordpress(w http.ResponseWriter, r *http.Request) {
	db := Db()
	defer db.Close()

	requestSchema := struct {
		UserId string `json:"user_id"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&requestSchema); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	userId := requestSchema.UserId
	username,_ := getUserNameByUserID(db, userId)

	if os.Getenv("GO_ENV") == "production" {
		// 実行するシェルスクリプトファイルのパス
		scriptPath := "../../script/wordpress_setup.sh"
	
		// シェルスクリプトの実行
		stdout, stderr, err := runShellScript(scriptPath, username)
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println("Stderr:", stderr)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	
		// 正常に終了した場合の処理
		fmt.Println("Script executed successfully")
		fmt.Println("Stdout:", stdout)
	}

	w.WriteHeader(http.StatusOK)
}

// GET localhost:3000/wordpress/status?userId={userID}
func getWordpressStatus(w http.ResponseWriter, r *http.Request) {
	db := Db()
	defer db.Close()

	userId := r.URL.Query().Get("userId")
	username,_ := getUserNameByUserID(db, userId)

	wordpressURL := fmt.Sprintf("https://crane74.com/%s/wordpress", username)
	
	statusCode, status, err := checkWebsiteStatus(wordpressURL)
    message := "Website is " + status

    if err != nil {
        message = err.Error()
    }

    response := StatusResponse{
        StatusCode: statusCode,
        Status:     status,
        Message:    message,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
