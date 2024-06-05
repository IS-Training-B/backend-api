package main

import (
	"encoding/json"
	"net/http"
)

type Mail struct {
	Id int `json:"id"`
	MailUserName string `json:"mail_username"`
	MailAddress string `json:"mail_address"`
}

type Mails []Mail

func getUserMails(w http.ResponseWriter, r *http.Request) {
	db := Db()
	defer db.Close()

	userId := r.URL.Query().Get("user_id")
	rows, err := db.Query("SELECT id, mail_username, mail_address FROM mails WHERE user_id=? AND deleted_at IS NOT NULL", userId)

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