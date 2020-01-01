package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	// dbUser := os.Getenv("DBUSER")
	// dbPass := os.Getenv("DBPASS")
	// dbname := os.Getenv("DBNAME")
	dbCreds := os.Getenv("CLEARDB_DATABASE_URL")
	// db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbname)
	db, err := sql.Open(dbDriver, dbCreds)
	if err != nil {
		fmt.Println("ERROR Connecting to DB")
		panic(err.Error())
	}
	fmt.Println("DB Success")
	return db
}

func main() {
	serveWeb()
}

// CheckAndUpdate checks open binance orders, updates their status, then updates dishes.
func CheckAndUpdate() {
	CheckAllOpenOrders()
	UpdateDishes()
}
