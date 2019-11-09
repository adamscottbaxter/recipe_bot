package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Dish struct {
	ID           int
	RecipeID     int
	PairOne      string
	PairTwo      string
	CurrentPrice float64
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "trade"
	dbPass := "tradebot"
	dbname := "trade_bot"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbname)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func main() {
	serveWeb()
}
