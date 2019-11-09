package main

import (
	"context"
	"fmt"
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

func CreateOrder(symbol string, quantity string, sellPrice string) *binance.CreateOrderResponse {
	order, err := CreateClient().NewCreateOrderService().Symbol(symbol).
		Side(binance.SideTypeSell).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(quantity).
		Price(sellPrice).Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("ORDER: \n %+v\n", order)
	return order
}

func CreateStopLossLimitOrder(symbol string, quantity string, sellPrice string, stopPrice string) *binance.CreateOrderResponse {
	order, err := CreateClient().NewCreateOrderService().Symbol(symbol).
		Side(binance.SideTypeSell).Type(binance.OrderTypeStopLossLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(quantity).
		Price(sellPrice).StopPrice(stopPrice).Do(context.Background())
	if err != nil {
		panic(err)
	}
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
