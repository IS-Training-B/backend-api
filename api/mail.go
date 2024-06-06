package main

import (
	"encoding/json"
	"database/sql"
	"net/http"
	"fmt"
	"os"
	"github.com/go-sql-driver/mysql"
)

type Mail struct {
	Id int `json:"id"`
	MailLocalpart string `json:"mail_localpart"`
	MailAddress string `json:"mail_address"`
}

type MailRequest struct {
    UserId string `json:"user_id"`
    MailLocalpart string `json:"mail_localpart"`
}

type Mails []Mail

// GET localhost:3000/mails?user_id={userID}
func getUserMails(w http.ResponseWriter, r *http.Request) {
	db := Db()
	defer db.Close()

	// Queryの取得
	userId := r.URL.Query().Get("userId")

	rows, err := db.Query("SELECT id, mail_localpart, mail_address FROM mails WHERE user_id=?", userId)

	var mails Mails
	 
	if err != nil {
		return
	}

	for rows.Next() {
		mail := Mail{}
		rows.Scan(&mail.Id, &mail.MailLocalpart, &mail.MailAddress)
		mails = append(mails, mail)
	}

	w.WriteHeader(http.StatusOK) // 200を返却
	json.NewEncoder(w).Encode(mails)
}

// POST localhost:3000/mail/create
func createMailAddress(w http.ResponseWriter, r *http.Request) {
	db := Db()
	defer db.Close()

	var requestSchema MailRequest

    if err := json.NewDecoder(r.Body).Decode(&requestSchema); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	userId := requestSchema.UserId
	localpart := requestSchema.MailLocalpart
	address := fmt.Sprintf("%s@%s", localpart, os.Getenv("DOMAIN"))
	username, err := getUserNameByUserID(db, userId); 
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
        return
	}

	// TODO: 動作確認（APIを叩いて正常にシェルスクリプトが走るか）
	// 新しいメールアドレスの追加処理

	if os.Getenv("GO_ENV") == "production" {
		// 実行するシェルスクリプトファイルのパス
		scriptPath := "../../script/mail_add.sh"
	
		// シェルスクリプトの実行
		stdout, stderr, err := runShellScript(scriptPath, username, localpart)
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
	

	sql, err := db.Prepare("INSERT INTO mails (user_id, mail_localpart, mail_address) VALUES (?,?,?)")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	_, err = sql.Exec(userId, localpart, address)
	if err != nil {
		// エラーがユニーク制約違反であるかをチェック
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			http.Error(w, "メールアドレスが既に存在します", http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	var (
		id           int
		MailLocalpart string
		mailAddress  string
	)

	// DBへの追加確認
	err = db.QueryRow("SELECT id, mail_localpart, mail_address FROM mails WHERE user_id=? ORDER BY created_at DESC LIMIT 1", userId).Scan(&id, &MailLocalpart, &mailAddress)
	if err != nil {
		fmt.Println("No rows found")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Mail{Id: id, MailLocalpart: MailLocalpart, MailAddress: mailAddress})
}

// POST localhost:3000/mail/delete
func deleteMailAddresss(w http.ResponseWriter, r *http.Request) {
	db := Db()
	defer db.Close()

	var requestSchema MailRequest

    if err := json.NewDecoder(r.Body).Decode(&requestSchema); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	userId := requestSchema.UserId
	MailLocalpart := requestSchema.MailLocalpart

	exists, err := checkRecordExists(db, MailLocalpart, userId)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if !exists {
        http.Error(w, "Record not found", http.StatusNotFound)
        return
    }

	// TODO: 動作確認（APIを叩いて正常にシェルスクリプトが走るか）
	// メールアドレスの削除処理

	// mail_delete.shが配置されたら下をコメントアウト

	// if os.Getenv("GO_ENV") == "production" {
	// 	// 実行するシェルスクリプトファイルのパス
	// 	scriptPath := "../../script/mail_delete.sh"
	//  http.Error(w, err.Error(), http.StatusBadRequest)
	// 	// シェルスクリプトの実行
	// 	stdout, stderr, err := runShellScript(scriptPath, username, localpart)
	// 	if err != nil {
	// 		fmt.Println("Error:", err)
	// 		fmt.Println("Stderr:", stderr)
	// 		http.Error(w, err.Error(), http.StatusBadRequest)
	// 		return
	// 	}
	
	// 	// 正常に終了した場合の処理
	// 	fmt.Println("Script executed successfully")
	// 	fmt.Println("Stdout:", stdout)
	// }

	if sql, err := db.Prepare("DELETE FROM mails WHERE mail_localpart = ? AND user_id = ?"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		if _, err := sql.Exec(MailLocalpart, userId); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func checkRecordExists(db *sql.DB, MailLocalpart string, userId string) (bool, error) {
    var exists bool
    query := "SELECT EXISTS(SELECT 1 FROM mails WHERE mail_localpart = ? AND user_id = ?)"
    err := db.QueryRow(query, MailLocalpart, userId).Scan(&exists)
    if err != nil {
        return false, err
    }
    return exists, nil
}