package main

import (
	"fmt"
	"log"

	binance "github.com/adshao/go-binance"
)

type Order struct {
	ID               int
	DishID           int
	Symbol           string
	BinanceOrderID   int
	BinanceStatus    string
	OriginalQuantity float64
	Price            float64
	NetChange        float64
}

func (o Order) UpdateStatus() string {
	fmt.Printf("ORDER update: \n %+v\n", o)
	binanceOrder := CheckOrder(o.Symbol, int64(o.BinanceOrderID))
	fmt.Printf("Binance Order Status: %T type = %v", binanceOrder.Status, binanceOrder.Status)
	o.setStatus(binanceOrder.Status)

	return "TBD"
}

func (o Order) setStatus(status binance.OrderStatusType) Order {
	db := dbConn()
	prep, err := db.Prepare("UPDATE orders SET binance_status=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	prep.Exec(status, o.ID)
	fmt.Println("UPDATE ORDER: status: %v | ID: %v", status, o.ID)
	defer db.Close()
	return o
}

func CheckAllOpenOrders() string {

	db := dbConn()

	rows, err := db.Query("SELECT id, dish_id, symbol, binance_order_id, binance_status, original_quantity, price, coalesce(net_change, 0.0) FROM orders WHERE binance_status = ?", "OPEN")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	order := Order{}
	orders := []Order{}

	var (
		id, dishID, binanceOrderID         int
		symbol, binanceStatus              string
		originalQuantity, price, netChange float64
	)
	for rows.Next() {
		err := rows.Scan(&id, &dishID, &symbol, &binanceOrderID, &binanceStatus, &originalQuantity, &price, &netChange)
		if err != nil {
			log.Fatal(err)
		}
		order.ID = id
		order.DishID = dishID
		order.Symbol = symbol
		order.BinanceOrderID = binanceOrderID
		order.BinanceStatus = binanceStatus
		order.OriginalQuantity = originalQuantity
		order.Price = price
		order.NetChange = netChange

		orders = append(orders, order)
	}

	for _, order := range orders {
		order.UpdateStatus()
	}

	return "TBD ALL OPEN ORDERS"
}