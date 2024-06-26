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
func getUserNameByUserID(db *sql.DB, userId string) (string, error) {
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

// DBの存在確認（Wordpress用）
func databaseExists(db *sql.DB, dbName string) (bool, error) {
    var exists bool
    query := fmt.Sprintf("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '%s'", dbName)
    err := db.QueryRow(query).Scan(&exists)
    if err == sql.ErrNoRows {
        return false, nil
    }
    if err != nil {
        return false, err
    }
    return true, nil
}

// ユーザIDからemail取得
func getUserEmail(db *sql.DB, userId string) (string, error) {
    var email string
    query := "SELECT email FROM users WHERE id = ?"

    err := db.QueryRow(query, userId).Scan(&email)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", fmt.Errorf("no user found with id %d")
        }
        return "", err
    }
    return email, nil
}