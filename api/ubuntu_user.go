package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type UserRequest struct {
    UserId string `json:"user_id"`
    UserName string `json:"username"`
}

func addUbuntuUser(w http.ResponseWriter, r *http.Request) {
	db := Db()
	defer db.Close()

	var requestSchema UserRequest
	if err := json.NewDecoder(r.Body).Decode(&requestSchema); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	// userId := requestSchema.UserId
	username := requestSchema.UserName

	if exist,err := checkUserNameExist(db, username); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if exist {
		http.Error(w, "そのユーザは既に存在しています", http.StatusBadRequest)
		return
	}

	// Ubuntuユーザのパスワード。user名と同じにしておく。
	password := username
	
	// TODO: 動作確認（APIを叩いて正常にシェルスクリプトが走るか）
	// 新規Ubuntuユーザーの追加処理
	
	if os.Getenv("GO_ENV") == "production" {
		// 実行するシェルスクリプトファイルのパス
		scriptPath := "../../script/user_add.sh"
	
		// シェルスクリプトの実行
		stdout, stderr, err := runShellScript(scriptPath, username, password)
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

func deleteUbuntuUser(w http.ResponseWriter, r *http.Request) {
	db := Db()
	defer db.Close()

	var requestSchema UserRequest
	if err := json.NewDecoder(r.Body).Decode(&requestSchema); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	// userId := requestSchema.UserId
	username := requestSchema.UserName

	// TODO: 動作確認（APIを叩いて正常にシェルスクリプトが走るか）
	// Ubuntuユーザーの削除処理

	if os.Getenv("GO_ENV") == "production" {
		// 実行するシェルスクリプトファイルのパス
		scriptPath := "../../script/delete_user.sh"
	
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