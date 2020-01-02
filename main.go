package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/adamscottbaxter/go-binance"
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
	fmt.Println("DB Success")
	return db
}

func main() {
	// serveWeb()
	symbolLevels := map[string]string{
		"BTCUSDT": "5",
		"ETHUSDT": "5",
	}

	// wsDepthHandler := func(event *binance.WsDepthEvent) {
	// 	fmt.Println(event)
	// }
	errHandler := func(err error) {
		fmt.Println(err)
	}

	wsPartial := func(event *binance.WsPartialDepthEvent) {
		fmt.Println(event)
	}

	// var wsPartialDepthEvt binance.WsPartialDepthEvent
	// var wsPartialDepthHandler wsPartial// binance.WsPartialDepthHandler
	// wsPartialDepthHandler := binance.WsPartialDepthHandler(func(event *wsPartialDepthEvt))

	doneC, stopC, err := binance.WsCombinedPartialDepthServe(symbolLevels, wsPartial, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	// use stopC to exit
	go func() {
		time.Sleep(3 * time.Second)
		stopC <- struct{}{}
	}()
	// remove this if you do not want to be blocked here
	<-doneC
}

// CheckAndUpdate checks open binance orders, updates their status, then updates dishes.
func CheckAndUpdate() {
	CheckAllOpenOrders()
	UpdateDishes()
}
