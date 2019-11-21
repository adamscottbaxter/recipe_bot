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
	dishID := r.prepDish()
	r.CreateOrders(dishID)
	// &{Symbol:BNBBTC OrderID:289673452 ClientOrderID:ywzi6PNQtfbeS5LeMFMUw2 TransactTime:1573336057195 Price:0.00225170 OrigQuantity:0.10000000 ExecutedQuantity:0.00000000 CummulativeQuoteQuantity:0.00000000 Status:NEW TimeInForce:GTC Type:LIMIT Side:SELL Fills:[]}

	// &{Symbol:BNBBTC OrderID:289673453 ClientOrderID:Gy7KR6i6dN08euRbRXXHoE TransactTime:1573336057397 Price: OrigQuantity: ExecutedQuantity: CummulativeQuoteQuantity: Status: TimeInForce: Type: Side: Fills:[]}
	return "Dish Cooked"
}

func (r Recipe) prepDish() int64 {
	currentPrice := GetPrice(r.Symbol)

	db := dbConn()

	stmt, err := db.Prepare("INSERT INTO dishes(recipe_id, symbol, side, current_price) VALUES(?,?,?,?)")
	if err != nil {
		log.Fatal("Cannot prepare DB statement", err)
	}
	res, err := stmt.Exec(r.ID, r.Symbol, r.Side, currentPrice)
	if err != nil {
		log.Fatal("Cannot run insert statement", err)
	}
	dishID, _ := res.LastInsertId()
	fmt.Printf("Inserted row: %d", dishID)
	log.Printf("INSERT DISH: id: %v | Symbol: %v | Side: %v | Current Price: %v", r.ID, r.Symbol, r.Side, currentPrice)

	defer db.Close()

	return dishID
}

func (r Recipe) binanceSide() binance.SideType {
	if r.Side == "SELL" {
		return binance.SideTypeSell
	} else {
		return binance.SideTypeSell
	}
}

func (r Recipe) CreateOrders(dishID int64) [2]*binance.CreateOrderResponse {
	currentPrice := GetPrice(r.Symbol)

	highPrice := MultiplyPrice(currentPrice, r.GainRatio)
	lowPrice := MultiplyPrice(currentPrice, r.LossRatio)
	stopPrice := MultiplyPrice(lowPrice, 0.99)
	highOrder := CreateOrder(dishID, r.Symbol, r.binanceSide(), r.StringQty(), highPrice)
	lowOrder := CreateStopLossLimitOrder(dishID, r.Symbol, r.binanceSide(), r.StringQty(), lowPrice, stopPrice)

	var orders [2]*binance.CreateOrderResponse
	orders[0] = highOrder
	orders[1] = lowOrder

	fmt.Println("Orders: %v", &orders)
	return orders
}

func (r Recipe) StringQty() string {
	return strconv.FormatFloat(r.Quantity, 'f', 7, 64)
}
