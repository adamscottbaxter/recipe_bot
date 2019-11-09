package main

import (
	"database/sql"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Recipe struct {
	ID        int
	Name      string
	Symbol    string
	Side      string
	GainRatio float64
	LossRatio float64
	Quantity  float64
	Frequency int
}

func (r Recipe) CreateOrders() string {
	currentPrice := GetPrice(r.Symbol)

	highPrice := MultiplyPrice(currentPrice, r.GainRatio)
	lowPrice := MultiplyPrice(currentPrice, r.LossRatio)
	stopPrice := MultiplyPrice(lowPrice, 0.99)
	CreateOrder(r.Symbol, r.StringQty(), highPrice)
	CreateStopLossLimitOrder(r.Symbol, r.StringQty(), lowPrice, stopPrice)
	return "ok"
}

func (r Recipe) StringQty() string {
	return strconv.FormatFloat(r.Quantity, 'f', 7, 64)
}

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
