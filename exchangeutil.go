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

func InsertValidOrder(dishID int64, symbol string, orderID int64, orderStatus binance.OrderStatusType, quantity string, price string) {
	db := dbConn()
	stmt, err := db.Prepare("INSERT INTO orders(dish_id, symbol, binance_order_id, binance_status, original_quantity, price) VALUES(?,?,?,?,?,?)")
	if err != nil {
		log.Fatal("Cannot prepare DB statement", err)
	}
	defer db.Close()
	_, err = stmt.Exec(dishID, symbol, orderID, orderStatus, quantity, price)
	if err != nil {
		panic(err)
		log.Fatal("Cannot run insert statement", err)
	}
	defer db.Close()
	fmt.Print("Created Limit Order")
}

func InsertInvalidOrder(dishID int64, symbol string, quantity string, price string, errorMessage string) {
	db := dbConn()
	stmt, err := db.Prepare("INSERT INTO orders(dish_id, symbol, original_quantity, price, error_message) VALUES(?,?,?,?,?)")
	if err != nil {
		log.Fatal("Cannot prepare DB statement", err)
	}
	_, err = stmt.Exec(dishID, symbol, quantity, price, errorMessage)
	if err != nil {
		panic(err)
		log.Fatal("Cannot run insert statement", err)
	}
	defer db.Close()
	fmt.Print("Insert INVALID order.")
}

func CreateOrder(dishID int64, symbol string, side binance.SideType, quantity string, sellPrice string) {
	// var binanceErr string
	order, err := CreateClient().NewCreateOrderService().Symbol(symbol).
		Side(side).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(quantity).
		Price(sellPrice).Do(context.Background())
	if err != nil {
		binanceErr := err.Error()
		InsertInvalidOrder(dishID, symbol, quantity, sellPrice, binanceErr)
	} else {
		InsertValidOrder(dishID, symbol, order.OrderID, order.Status, order.OrigQuantity, order.Price)
	}
	// fmt.Printf("ORDER: \n %+v\n error- %v", order, binanceErr)
	// db := dbConn()
	// stmt, err := db.Prepare("INSERT INTO orders(dish_id, symbol) VALUES(?,?)")
	// if err != nil {
	// 	log.Fatal("Cannot prepare DB statement", err)
	// }
	// fmt.Print("Created Limit Order")

	// &{Symbol:BNBBTC OrderID:289673452 ClientOrderID:ywzi6PNQtfbeS5LeMFMUw2 TransactTime:1573336057195 Price:0.00225170 OrigQuantity:0.10000000 ExecutedQuantity:0.00000000 CummulativeQuoteQuantity:0.00000000 Status:NEW TimeInForce:GTC Type:LIMIT Side:SELL Fills:[]}
	// _, err = stmt.Exec(dishID, order.Symbol, order.OrderID, order.Status, order.OrigQuantity, order.Price, binanceErr)
	// // stmt.Exec(dishID, order.Symbol, order.OrderID, order.Status, order.OrigQuantity, order.Price, binanceErr)
	// if err != nil {
	// 	panic(err)
	// 	log.Fatal("Cannot run insert statement", err)
	// 	return order
	// }
	// defer db.Close()

}

func CreateOrderTest(symbol string, side binance.SideType, quantity string, sellPrice string) error {
	err := CreateClient().NewCreateOrderService().Symbol(symbol).
		Side(side).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(quantity).
		Price(sellPrice).Test(context.Background())

	return err
}

func CreateStopLossLimitOrder(dishID int64, symbol string, side binance.SideType, quantity string, sellPrice string, stopPrice string) {
	order, err := CreateClient().NewCreateOrderService().Symbol(symbol).
		Side(side).Type(binance.OrderTypeStopLossLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(quantity).
		Price(sellPrice).StopPrice(stopPrice).Do(context.Background())

	if err != nil {
		binanceErr := err.Error()
		InsertInvalidOrder(dishID, symbol, quantity, sellPrice, binanceErr)
	} else {
		InsertValidOrder(dishID, symbol, order.OrderID, order.Status, order.OrigQuantity, order.Price)
	}

	// if err != nil {
	// 	panic(err)
	// }
	// db := dbConn()
	// stmt, err := db.Prepare("INSERT INTO orders(dish_id, symbol, binance_order_id, binance_status, original_quantity, price) VALUES(?,?,?,?,?,?)")
	// if err != nil {
	// 	log.Fatal("Cannot prepare DB statement", err)
	// }
	// fmt.Print("Created Stop Loss Order")
	// // &{Symbol:BNBBTC OrderID:289673453 ClientOrderID:Gy7KR6i6dN08euRbRXXHoE TransactTime:1573336057397 Price: OrigQuantity: ExecutedQuantity: CummulativeQuoteQuantity: Status: TimeInForce: Type: Side: Fills:[]}
	// stmt.Exec(dishID, order.Symbol, order.OrderID, order.Status, order.OrigQuantity, order.Price)
	// if err != nil {
	// 	log.Fatal("Cannot run insert statement", err)
	// }
	// defer db.Close()
	// fmt.Printf("ORDER: \n %+v\n", order)
	// return order
}

func CreateStopLossLimitOrderTest(symbol string, side binance.SideType, quantity string, sellPrice string, stopPrice string) error {
	err := CreateClient().NewCreateOrderService().Symbol(symbol).
		Side(side).Type(binance.OrderTypeStopLossLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(quantity).
		Price(sellPrice).StopPrice(stopPrice).Test(context.Background())

	return err
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

func CreateTakeProfitLimitOrderTest(symbol string, side binance.SideType, quantity string, sellPrice string, stopPrice string) error {
	err := CreateClient().NewCreateOrderService().Symbol(symbol).
		Side(side).Type(binance.OrderTypeTakeProfitLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(quantity).
		Price(sellPrice).StopPrice(stopPrice).Test(context.Background())

	return err
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
	fmt.Printf("ORDER CHECK: \n %+v\n", order)
	return order
}

// https://www.binance.com/api/v1/exchangeInfo
// {
//   "symbol": "BNBBTC",
//   "status": "TRADING",
//   "baseAsset": "BNB",
//   "baseAssetPrecision": 8,
//   "quoteAsset": "BTC",
//   "quotePrecision": 8,
//   "baseCommissionPrecision": 8,
//   "quoteCommissionPrecision": 8,
//   "orderTypes": [
//     "LIMIT",
//     "LIMIT_MAKER",
//     "MARKET",
//     "STOP_LOSS_LIMIT",
//     "TAKE_PROFIT_LIMIT"
//   ],
//   "icebergAllowed": true,
//   "ocoAllowed": true,
//   "quoteOrderQtyMarketAllowed": true,
//   "isSpotTradingAllowed": true,
//   "isMarginTradingAllowed": true,
//   "filters": [
//     {
//       "filterType": "PRICE_FILTER",
//       "minPrice": "0.00000010",
//       "maxPrice": "100000.00000000",
//       "tickSize": "0.00000010"
//     },
//     {
//       "filterType": "PERCENT_PRICE",
//       "multiplierUp": "5",
//       "multiplierDown": "0.2",
//       "avgPriceMins": 5
//     },
//     {
//       "filterType": "LOT_SIZE",
//       "minQty": "0.01000000",
//       "maxQty": "100000.00000000",
//       "stepSize": "0.01000000"
//     },
//     {
//       "filterType": "MIN_NOTIONAL",
//       "minNotional": "0.00010000",
//       "applyToMarket": true,
//       "avgPriceMins": 5
//     },
//     {
//       "filterType": "ICEBERG_PARTS",
//       "limit": 10
//     },
//     {
//       "filterType": "MARKET_LOT_SIZE",
//       "minQty": "0.00000000",
//       "maxQty": "1769700.00000000",
//       "stepSize": "0.00000000"
//     },
//     {
//       "filterType": "MAX_NUM_ALGO_ORDERS",
//       "maxNumAlgoOrders": 5
//     }
//   ]
// }
