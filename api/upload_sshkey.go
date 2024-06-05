package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type SSHkeyRequest struct {
    UserId int `json:"user_id"`
    SSHkey string `json:"ssh_public_key"`
}


// TODO: 動作確認（APIを叩いて正常にシェルスクリプトが走るか）
// SSH公開鍵のアップロード

// POST localhost:3000/sshkey/upload
func uploadSSHkey(w http.ResponseWriter, r *http.Request) {
	db := Db()
	defer db.Close()

	var requestSchema SSHkeyRequest

    if err := json.NewDecoder(r.Body).Decode(&requestSchema); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	userId := requestSchema.UserId
	sshkey := requestSchema.SSHkey
	username, err := getUserNameByUserID(db, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
        return
	}

	if os.Getenv("GO_ENV") == "production" {
		// 実行するシェルスクリプトファイルのパス
		scriptPath := "../../script/ssh_key_add.sh"
	
		// シェルスクリプトの実行
		stdout, stderr, err := runShellScript(scriptPath, username, sshkey)
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