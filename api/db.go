package main

import (
	"database/sql"
	"fmt"
    "os"
	_ "github.com/go-sql-driver/mysql"
)

func Db() (db *sql.DB)  {
	// db, _ := sql.Open("mysql", "[username]:[password]@/[database]")
    db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", os.Getenv("USER_NAME"), os.Getenv("PASSWORD"), os.Getenv("DATABASE_NAME")))
    if err != nil {
		panic(err.Error())
	}
    return db
}

// ユーザIDからユーザ名を取得
func getUserNameByUserID(db *sql.DB, userId int) (string, error) {
    var username string
    query := "SELECT name FROM users WHERE id = ?"

    err := db.QueryRow(query, userId).Scan(&username)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", fmt.Errorf("no user found with id %d", userId)
        }
        return "", err
    }
    return username, nil
}

// ユーザ名の存在確認
func checkUserNameExist(db *sql.DB, name string) (bool, error) {
    var exist bool
    query := "SELECT EXISTS(SELECT 1 FROM users WHERE name = ?)"

    err := db.QueryRow(query, name).Scan(&exist)
    if err != nil {
        return false, err
    }
    return exist, nil
}