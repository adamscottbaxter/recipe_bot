package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

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
	// CheckAllOpenOrders()
	// db := dbConn()
	// nId := 8
	// selDB, err := db.Query("SELECT * FROM recipes WHERE id=?", nId)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// recipe := Recipe{}
	// for selDB.Next() {
	// 	var id int
	// 	var name, symbol, side string
	// 	var gainRatio, lossRatio, quantity float64
	// 	var frequency int
	// 	err = selDB.Scan(&id, &name, &symbol, &side, &gainRatio, &lossRatio, &quantity, &frequency)
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	// 	recipe.ID = id
	// 	recipe.Name = name
	// 	recipe.Symbol = symbol
	// 	recipe.Side = side
	// 	recipe.GainRatio = gainRatio
	// 	recipe.LossRatio = lossRatio
	// 	recipe.Quantity = quantity
	// 	recipe.Frequency = frequency
	// }

	// defer db.Close()

	// recipe.CookDish()
}
