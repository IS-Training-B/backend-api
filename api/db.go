package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func Db() {
    // Create the database handle, confirm driver is present
	// db, _ := sql.Open("mysql", "[username]:[password]@/[database]")
    db, _ := sql.Open("mysql", "root:password@/rs")
    defer db.Close()

    // Connect and check the server version
    var version string
    db.QueryRow("SELECT VERSION()").Scan(&version)
    fmt.Println("Connected to: ", version)


}