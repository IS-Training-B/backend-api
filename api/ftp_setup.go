package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// TODO: 動作確認（APIを叩いて正常にシェルスクリプトが走るか）
// FTPの初期設定

// POST localhost:3000/ftp/setup
func setupFTP(w http.ResponseWriter, r *http.Request) {
	requestSchema := struct {
		UserId int `json:"user_id"`
        UserName string `json:"username"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&requestSchema); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	username := requestSchema.UserName

	if os.Getenv("GO_ENV") == "production" {
		// 実行するシェルスクリプトファイルのパス
		scriptPath := "../../script/ftp_access_control.sh"
	
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