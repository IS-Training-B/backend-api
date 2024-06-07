package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"database/sql"
)

type UserCreateRequest struct {
    UserId string `json:"user_id"`
    UserName string `json:"username"`
	Password string `json:"password"`
}

type UserDeleteRequest struct {
	UserId string `json:"user_id"`
    UserName string `json:"username"`
}

func addUbuntuUser(w http.ResponseWriter, r *http.Request) {
	db := Db()
	defer db.Close()

	var requestSchema UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&requestSchema); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	userId := requestSchema.UserId
	username := requestSchema.UserName
	password := requestSchema.Password

	if exist,err := checkUserNameExist(db, username); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if !exist {
		http.Error(w, "そのユーザは登録されていません", http.StatusBadRequest)
		return
	}

	// passwordをDBに登録
	err := setUserPassword(db, userId, password)
	if err != nil {
		http.Error(w, "DBへのパスワード登録に失敗しました", http.StatusBadRequest)
		return
	}
	
	// TODO: 動作確認（APIを叩いて正常にシェルスクリプトが走るか）
	// 新規Ubuntuユーザーの追加処理
	
	if os.Getenv("GO_ENV") == "production" {
		// 実行するシェルスクリプトファイルのパス
		scriptPath := "../../script/user_add.sh"
	
		// シェルスクリプトの実行
		fmt.Println(fmt.Sprintf("User: %s add script run ...", username))
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

	var requestSchema UserDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&requestSchema); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	userId := requestSchema.UserId
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

	// ユーザが作成したメールアドレスを全て削除
	deleteUserAllMail(userId)

	// ユーザのDB削除
	query := "DROP DATABASE IF EXISTS ?"
	stmt, err := db.Prepare(query)
    if err != nil {
        return
    }
    defer stmt.Close()
    _, err = stmt.Exec(username)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func setUserPassword(db *sql.DB, userID string, newPassword string) error {
    query := "UPDATE users SET password = ? WHERE id = ?"
    stmt, err := db.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(newPassword, userID)
    return err
}