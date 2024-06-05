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
	MailUserName string `json:"mail_localpart"`
	MailAddress string `json:"mail_address"`
}

type MailRequest struct {
    UserId int `json:"user_id"`
    MailUserName string `json:"mail_localpart"`
}

type Mails []Mail

func getUserMails(w http.ResponseWriter, r *http.Request) {
	db := Db()
	defer db.Close()

	// Queryの取得
	userId := r.URL.Query().Get("user_id")

	rows, err := db.Query("SELECT id, mail_localpart, mail_address FROM mails WHERE user_id=?", userId)

	var mails Mails
	if err != nil {
		return
	}

	for rows.Next() {
		mail := Mail{}
		rows.Scan(&mail.Id, &mail.MailUserName, &mail.MailAddress)
		mails = append(mails, mail)
	}

	w.WriteHeader(http.StatusOK) // 200を返却
	json.NewEncoder(w).Encode(mails)
}

func createMailAddress(w http.ResponseWriter, r *http.Request) {
	db := Db()
	defer db.Close()

	var requestSchema MailRequest

    if err := json.NewDecoder(r.Body).Decode(&requestSchema); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }


	// TODO: 新しいメールアドレスの追加処理


	userId := requestSchema.UserId
	username := requestSchema.MailUserName
	address := fmt.Sprintf("%s@%s", username, os.Getenv("DOMAIN"))
	
	sql, err := db.Prepare("INSERT INTO mails (user_id, mail_localpart, mail_address) VALUES (?,?,?)")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	_, err = sql.Exec(userId, username, address)
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
		mailUsername string
		mailAddress  string
	)

	// 追加確認
	err = db.QueryRow("SELECT id, mail_localpart, mail_address FROM mails WHERE user_id=? ORDER BY created_at DESC LIMIT 1", userId).Scan(&id, &mailUsername, &mailAddress)
	if err != nil {
		fmt.Println("No rows found")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Mail{Id: id, MailUserName: mailUsername, MailAddress: mailAddress})
}

func deleteMailAddresss(w http.ResponseWriter, r *http.Request) {
	db := Db()
	defer db.Close()

	var requestSchema MailRequest

    if err := json.NewDecoder(r.Body).Decode(&requestSchema); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	userId := requestSchema.UserId
	mailUsername := requestSchema.MailUserName

	exists, err := checkRecordExists(db, mailUsername, userId)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if !exists {
        http.Error(w, "Record not found", http.StatusNotFound)
        return
    }


	// TODO: メールアドレスの削除処理



	if sql, err := db.Prepare("DELETE FROM mails WHERE mail_localpart = ? AND user_id = ?"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		if _, err := sql.Exec(mailUsername, userId); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func checkRecordExists(db *sql.DB, mailUsername string, userId int) (bool, error) {
    var exists bool
    query := "SELECT EXISTS(SELECT 1 FROM mails WHERE mail_localpart = ? AND user_id = ?)"
    err := db.QueryRow(query, mailUsername, userId).Scan(&exists)
    if err != nil {
        return false, err
    }
    return exists, nil
}