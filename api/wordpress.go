package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"database/sql"
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
	ubuntuPassword, _ := getUserPassword(db, userId)
	email, _ := getUserEmail(db, userId)

	// wordpress用のユーザDBが存在するか確認
	exist,err := databaseExists(db, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
        return
	}

	// wordpress用のユーザDBがなければ作成
	if !exist && os.Getenv("GO_ENV") == "production" {
		// 実行するシェルスクリプトファイルのパス
		scriptPath := "../../script/make-user.sh"
	
		// シェルスクリプトの実行
		// $0 <db_user> <db_password>
		stdout, stderr, err := runShellScript(scriptPath, username, ubuntuPassword)
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println("Stderr:", stderr)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	
		// 正常に終了した場合の処理
		fmt.Println("create user's DB script executed successfully")
		fmt.Println("Stdout:", stdout)
	}

	if os.Getenv("GO_ENV") == "production" {
		// 実行するシェルスクリプトファイルのパス
		scriptPath := "../../script/wordpress_setup.sh"
	
		// シェルスクリプトの実行
		//  $0 <username> <domain> <db_password> <admin_password> <admin_email>
		stdout, stderr, err := runShellScript(scriptPath, username, os.Getenv("DOMAIN"), ubuntuPassword, ubuntuPassword, email)
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
	
	fmt.Println(fmt.Sprintf("URL: %s status get", wordpressURL))
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

func getUserPassword(db *sql.DB, userID string) (string, error) {
    var password string
    query := "SELECT password FROM users WHERE id = ?"
    err := db.QueryRow(query, userID).Scan(&password)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", fmt.Errorf("no user with id %s", userID)
        }
        return "", err
    }
    return password, nil
}