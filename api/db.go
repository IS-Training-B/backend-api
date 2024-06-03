package main

import (
	"database/sql"
	"fmt"
    "os"
	_ "github.com/go-sql-driver/mysql"
)

func Db() {
	// db, _ := sql.Open("mysql", "[username]:[password]@/[database]")
    db, _ := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", os.Getenv("USER_NAME"), os.Getenv("PASSWORD"), os.Getenv("DATABASE_NAME")))
    defer db.Close()

    // Connect and check the server version
    var version string
    db.QueryRow("SELECT VERSION()").Scan(&version)
    fmt.Println("Connected to: ", version)


}