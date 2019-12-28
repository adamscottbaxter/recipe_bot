package main

import (
	"fmt"
	"log"

	binance "github.com/adshao/go-binance"
)

// Order type
type Order struct {
	ID               int
	DishID           int
	Symbol           string
	BinanceOrderID   int
	BinanceStatus    string
	OriginalQuantity float64
	Price            float64
	ErrorMessage     string
}

// UpdateStatus update the status of an order
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
	fmt.Printf("UPDATE ORDER: status: %v | ID: %v", status, o.ID)
	defer db.Close()
	return o
}

// GetAllOrders gets all Orders without errors from db
func GetAllOrders() []Order {
	db := dbConn()

	rows, err := db.Query("SELECT id, dish_id, COALESCE(symbol, ''), COALESCE(binance_order_id, 0), COALESCE(binance_status, ''), COALESCE(original_quantity, 0), COALESCE(price, 0), COALESCE(error_message, '') FROM orders where error_message IS NULL")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	order := Order{}
	orders := []Order{}

	var (
		id, dishID, binanceOrderID          int
		symbol, binanceStatus, errorMessage string
		originalQuantity, price             float64
	)
	for rows.Next() {
		err := rows.Scan(&id, &dishID, &symbol, &binanceOrderID, &binanceStatus, &originalQuantity, &price, &errorMessage)
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
		order.ErrorMessage = errorMessage

		orders = append(orders, order)
	}
	fmt.Printf("GetAllOrders: %v", orders)
	return orders
}

// CheckAllOpenOrders updates the status of all open orders
func CheckAllOpenOrders() string {
	orders := GetAllOrders()

	for _, order := range orders {
		if order.BinanceStatus == "NEW" {
			order.UpdateStatus()
		}
	}

	return "TBD ALL OPEN ORDERS"
}
