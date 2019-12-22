package main

import "fmt"

type Dish struct {
	ID           int
	RecipeID     int
	Symbol       string
	Side         string
	CurrentPrice float64
	FillPrice    float64
	NetChange    float64
}

func AllDishes() []Dish {
	db := dbConn()
	selDB, err := db.Query("SELECT id, recipe_id, symbol, side, COALESCE(current_price, 0), COALESCE(fill_price, 0), COALESCE(net_change, 0) FROM dishes ORDER BY id ASC")
	if err != nil {
		panic(err.Error())
	}
	dish := Dish{}
	dishes := []Dish{}
	for selDB.Next() {
		var (
			id, recipeID                       int
			symbol, side                       string
			currentPrice, fillPrice, netChange float64
		)

		err = selDB.Scan(&id, &recipeID, &symbol, &side, &currentPrice, &fillPrice, &netChange)
		if err != nil {
			panic(err.Error())
		}
		dish.ID = id
		dish.RecipeID = recipeID
		dish.Symbol = symbol
		dish.Side = side
		dish.CurrentPrice = currentPrice
		dish.FillPrice = fillPrice
		dish.NetChange = netChange
		dishes = append(dishes, dish)
	}
	defer db.Close()
	fmt.Print("Dishes: %v", dishes)
	return dishes
}

func (d Dish) Orders() []Order {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM orders where dish_id = ?", d.ID)
	if err != nil {
		panic(err.Error())
	}
	order := Order{}
	orders := []Order{}
	for selDB.Next() {
		var (
			id, dishID, binanceOrderID          int
			symbol, binanceStatus, errorMessage string
			originalQuantity, price             float64
		)

		err = selDB.Scan(&id, &dishID, &symbol, &binanceOrderID, &binanceStatus, &originalQuantity, &price, &errorMessage)
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
		order.ErrorMessage = errorMessage
		orders = append(orders, order)
	}
	defer db.Close()
	fmt.Print("Dish: %v ---- Orders: %v", d, orders)
	return orders
}
