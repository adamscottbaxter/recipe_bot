package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbname := os.Getenv("DBNAME")
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbname)
	if err != nil {
		fmt.Println("ERROR Connecting to DB")
		panic(err.Error())
	}
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
