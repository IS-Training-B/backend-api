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