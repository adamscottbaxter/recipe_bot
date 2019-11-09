package main

import (
	"fmt"
	"log"
	"strconv"

	binance "github.com/adshao/go-binance"
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

func (r Recipe) CookDish() string {
	currentPrice := GetPrice(r.Symbol)

	db := dbConn()

	insert, err := db.Prepare("INSERT INTO dishes(recipe_id, symbol, side, current_price) VALUES(?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	insert.Exec(r.ID, r.Symbol, r.Side, currentPrice)
	log.Printf("INSERT DISH: id: %v | Symbol: %v | Side: %v | Current Price: %v", r.ID, r.Symbol, r.Side, currentPrice)

	defer db.Close()

	return "Dish Cooked"
}

func (r Recipe) CreateOrders() [2]*binance.CreateOrderResponse {
	currentPrice := GetPrice(r.Symbol)

	highPrice := MultiplyPrice(currentPrice, r.GainRatio)
	lowPrice := MultiplyPrice(currentPrice, r.LossRatio)
	stopPrice := MultiplyPrice(lowPrice, 0.99)
	highOrder := CreateOrder(r.Symbol, r.StringQty(), highPrice)
	lowOrder := CreateStopLossLimitOrder(r.Symbol, r.StringQty(), lowPrice, stopPrice)

	var orders [2]*binance.CreateOrderResponse
	orders[0] = highOrder
	orders[1] = lowOrder

	fmt.Println("Orders: %v", orders)
	return orders
}

func (r Recipe) StringQty() string {
	return strconv.FormatFloat(r.Quantity, 'f', 7, 64)
}
