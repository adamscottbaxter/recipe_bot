package main

import (
	"fmt"

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

	selDB, err := db.Query("SELECT id, dish_id, symbol, binance_order_id, binance_status, original_quantity, price, coalesce(net_change, 0.0) FROM orders WHERE binance_status <> ?", "FILLED")
	if err != nil {
		panic(err.Error())
	}
	order := Order{}
	orders := []Order{}
	for selDB.Next() {
		var id, dishID, binanceOrderID int
		var symbol, binanceStatus string
		var originalQuantity, price, netChange float64
		err = selDB.Scan(&id, &dishID, &symbol, &binanceOrderID, &binanceStatus, &originalQuantity, &price, &netChange)
		if err != nil {
			panic(err.Error())
		}
		order.ID = id
		order.DishID = dishID
		order.Symbol = symbol
		order.BinanceOrderID = binanceOrderID
		order.BinanceStatus = binanceStatus
		order.OriginalQuantity = originalQuantity
		order.Price = price
		order.NetChange = netChange
	}

	orders = append(orders, order)

	for _, order := range orders {
		order.UpdateStatus()
	}
	// tmpl.ExecuteTemplate(w, "Show", recipe)
	defer db.Close()

	return "TBD ALL OPEN ORDERS"
}
