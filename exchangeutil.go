package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	binance "github.com/adshao/go-binance"
)

func CreateClient() *binance.Client {
	var (
		apiKey    = os.Getenv("GOBOTAPIKEY")
		secretKey = os.Getenv("GOBOTAPISECRET")
	)
	client := binance.NewClient(apiKey, secretKey)
	return client
}

func GetPrice(symbol string) string {
	prices, err := CreateClient().NewListPricesService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return "ERROR"
	}

	for _, p := range prices {
		if p.Symbol == symbol {
			return p.Price
		}
	}
	return "ERROR"
}

func MultiplyPrice(price string, factor float64) string {
	floatPrice, _ := strconv.ParseFloat(price, 64)
	floatPrice *= factor
	return strconv.FormatFloat(floatPrice, 'f', 7, 64)
}

func CreateOrder(dishID int64, symbol string, side binance.SideType, quantity string, sellPrice string) *binance.CreateOrderResponse {
	order, err := CreateClient().NewCreateOrderService().Symbol(symbol).
		Side(side).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(quantity).
		Price(sellPrice).Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("ORDER: \n %+v\n", order)
	db := dbConn()
	stmt, err := db.Prepare("INSERT INTO orders(dish_id, symbol, binance_order_id, binance_status, original_quantity, price) VALUES(?,?,?,?,?,?)")
	if err != nil {
		log.Fatal("Cannot prepare DB statement", err)
	}
	fmt.Print("Created Limit Order")

	// &{Symbol:BNBBTC OrderID:289673452 ClientOrderID:ywzi6PNQtfbeS5LeMFMUw2 TransactTime:1573336057195 Price:0.00225170 OrigQuantity:0.10000000 ExecutedQuantity:0.00000000 CummulativeQuoteQuantity:0.00000000 Status:NEW TimeInForce:GTC Type:LIMIT Side:SELL Fills:[]}
	stmt.Exec(dishID, order.Symbol, order.OrderID, order.Status, order.OrigQuantity, order.Price)
	if err != nil {
		log.Fatal("Cannot run insert statement", err)
	}
	defer db.Close()
	return order
}

func CreateStopLossLimitOrder(dishID int64, symbol string, side binance.SideType, quantity string, sellPrice string, stopPrice string) *binance.CreateOrderResponse {
	order, err := CreateClient().NewCreateOrderService().Symbol(symbol).
		Side(side).Type(binance.OrderTypeStopLossLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(quantity).
		Price(sellPrice).StopPrice(stopPrice).Do(context.Background())
	if err != nil {
		panic(err)
	}
	db := dbConn()
	stmt, err := db.Prepare("INSERT INTO orders(dish_id, symbol, binance_order_id, binance_status, original_quantity, price) VALUES(?,?,?,?,?,?)")
	if err != nil {
		log.Fatal("Cannot prepare DB statement", err)
	}
	fmt.Print("Created Stop Loss Order")
	// &{Symbol:BNBBTC OrderID:289673453 ClientOrderID:Gy7KR6i6dN08euRbRXXHoE TransactTime:1573336057397 Price: OrigQuantity: ExecutedQuantity: CummulativeQuoteQuantity: Status: TimeInForce: Type: Side: Fills:[]}
	stmt.Exec(dishID, order.Symbol, order.OrderID, order.Status, order.OrigQuantity, order.Price)
	if err != nil {
		log.Fatal("Cannot run insert statement", err)
	}
	defer db.Close()
	fmt.Printf("ORDER: \n %+v\n", order)
	return order
}

func CreateTakeProfitLimitOrder(dishID int64, symbol string, side binance.SideType, quantity string, sellPrice string, stopPrice string) *binance.CreateOrderResponse {
	order, err := CreateClient().NewCreateOrderService().Symbol(symbol).
		Side(side).Type(binance.OrderTypeTakeProfitLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(quantity).
		Price(sellPrice).StopPrice(stopPrice).Do(context.Background())
	if err != nil {
		panic(err)
	}
	db := dbConn()
	stmt, err := db.Prepare("INSERT INTO orders(dish_id, symbol, binance_order_id, binance_status, original_quantity, price) VALUES(?,?,?,?,?,?)")
	if err != nil {
		log.Fatal("Cannot prepare DB statement", err)
	}
	fmt.Print("Created Take Profit Limit Order")

	stmt.Exec(dishID, order.Symbol, order.OrderID, order.Status, order.OrigQuantity, order.Price)
	if err != nil {
		log.Fatal("Cannot run insert statement", err)
	}
	defer db.Close()
	fmt.Printf("ORDER: \n %+v\n", order)
	return order
}

func CancelOrder(symbol string, orderID int64) *binance.CancelOrderResponse {
	cancelledOrder, err := CreateClient().NewCancelOrderService().Symbol(symbol).
		OrderID(orderID).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return cancelledOrder
}

func CheckOrder(symbol string, orderID int64) *binance.Order {
	order, err := CreateClient().NewGetOrderService().Symbol(symbol).
		OrderID(orderID).Do(context.Background())
	if err != nil {
		panic(err)
	}
	return order
}
