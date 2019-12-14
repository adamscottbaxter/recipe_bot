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
	Active    bool
}

func CookActiveRecipes() {
	for _, r := range getActiveRecipes() {
		r.CookDish()
	}
}

func getActiveRecipes() []Recipe {
	db := dbConn()
	rows, err := db.Query("SELECT * FROM recipes where active = 1 ORDER BY id ASC")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	recipe := Recipe{}
	recipes := []Recipe{}
	for rows.Next() {
		var (
			id                             int
			name, symbol, side             string
			gainRatio, lossRatio, quantity float64
			frequency                      int
			active                         bool
		)

		err = rows.Scan(&id, &name, &symbol, &side, &gainRatio, &lossRatio, &quantity, &frequency, &active)
		if err != nil {
			panic(err.Error())
		}
		recipe.ID = id
		recipe.Name = name
		recipe.Symbol = symbol
		recipe.Side = side
		recipe.GainRatio = gainRatio
		recipe.LossRatio = lossRatio
		recipe.Quantity = quantity
		recipe.Frequency = frequency
		recipe.Active = active
		recipes = append(recipes, recipe)
	}
	fmt.Println("Recipes:", recipes)

	return recipes
}

func (r Recipe) CookDish() string {
	dishID := r.prepDish()
	if r.Side == "SELL" {
		r.CreateSellOrders(dishID)
	} else {
		r.CreateBuyOrders(dishID)
	}
	// &{Symbol:BNBBTC OrderID:289673452 ClientOrderID:ywzi6PNQtfbeS5LeMFMUw2 TransactTime:1573336057195 Price:0.00225170 OrigQuantity:0.10000000 ExecutedQuantity:0.00000000 CummulativeQuoteQuantity:0.00000000 Status:NEW TimeInForce:GTC Type:LIMIT Side:SELL Fills:[]}

	// &{Symbol:BNBBTC OrderID:289673453 ClientOrderID:Gy7KR6i6dN08euRbRXXHoE TransactTime:1573336057397 Price: OrigQuantity: ExecutedQuantity: CummulativeQuoteQuantity: Status: TimeInForce: Type: Side: Fills:[]}

	return "Dish Cooked"
}

func (r Recipe) CookDishTest() string {
	if r.Side == "SELL" {
		r.CreateSellOrdersTest()
	} else {
		r.CreateBuyOrdersTest()
	}
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
		return binance.SideTypeBuy
	}
}

func (r Recipe) CreateSellOrders(dishID int64) [2]*binance.CreateOrderResponse {
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

func (r Recipe) CreateSellOrdersTest() []error {
	currentPrice := GetPrice(r.Symbol)

	highPrice := MultiplyPrice(currentPrice, r.GainRatio)
	lowPrice := MultiplyPrice(currentPrice, r.LossRatio)
	stopPrice := MultiplyPrice(lowPrice, 0.99)
	highOrderError := CreateOrderTest(r.Symbol, r.binanceSide(), r.StringQty(), highPrice)
	lowOrderError := CreateStopLossLimitOrderTest(r.Symbol, r.binanceSide(), r.StringQty(), lowPrice, stopPrice)

	var orderErrors []error
	orderErrors = append(orderErrors, highOrderError)
	orderErrors = append(orderErrors, lowOrderError)

	fmt.Println("Order errors: %v", &orderErrors)
	return orderErrors
}

func (r Recipe) CreateBuyOrders(dishID int64) [2]*binance.CreateOrderResponse {
	currentPrice := GetPrice(r.Symbol)

	highPrice := MultiplyPrice(currentPrice, r.GainRatio)
	lowPrice := MultiplyPrice(currentPrice, r.LossRatio)
	stopPrice := MultiplyPrice(highPrice, 0.99)
	// see what high price and stop price are relative to current price
	lowOrder := CreateOrder(dishID, r.Symbol, r.binanceSide(), r.StringQty(), lowPrice)
	highOrder := CreateTakeProfitLimitOrder(dishID, r.Symbol, r.binanceSide(), r.StringQty(), highPrice, stopPrice)

	var orders [2]*binance.CreateOrderResponse
	orders[0] = highOrder
	orders[1] = lowOrder

	fmt.Println("Orders: %v", &orders)
	return orders
}

func (r Recipe) CreateBuyOrdersTest() []error {
	currentPrice := GetPrice(r.Symbol)

	highPrice := MultiplyPrice(currentPrice, r.GainRatio)
	lowPrice := MultiplyPrice(currentPrice, r.LossRatio)
	stopPrice := MultiplyPrice(highPrice, 0.999)

	lowOrderError := CreateOrderTest(r.Symbol, r.binanceSide(), r.StringQty(), lowPrice)
	highOrderError := CreateTakeProfitLimitOrderTest(r.Symbol, r.binanceSide(), r.StringQty(), highPrice, stopPrice)

	var orderErrors []error
	orderErrors = append(orderErrors, highOrderError)
	orderErrors = append(orderErrors, lowOrderError)

	fmt.Println("Order errors: %v", &orderErrors)
	return orderErrors
}

func (r Recipe) StringQty() string {
	return strconv.FormatFloat(r.Quantity, 'f', 7, 64)
}
