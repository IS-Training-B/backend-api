package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"os"
)

func example(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println(address)

	username, err := getUserNameByUserID(db, userId); 
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
        return
	}

	if os.Getenv("GO_ENV") == "production" {
		// 実行するシェルスクリプトファイルのパス
		scriptPath := "../../script/example.sh"
	
		// シェルスクリプトの実行
		stdout, stderr, err := runShellScript(scriptPath, username, localpart)
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println("Stderr:", stderr)
			fmt.Println("Stdout:", stdout)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println("Stdout:", stdout)
	}
}