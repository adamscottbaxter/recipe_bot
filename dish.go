package main

import (
	"fmt"
	"log"
	"math"
)

// Dish create dish
type Dish struct {
	ID           int
	RecipeID     int
	Symbol       string
	Side         string
	CurrentPrice float64
	FillPrice    float64
	NetChange    float64
}

// AllDishes returns a slice of all Dishes
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
	fmt.Printf("Dishes: %v", dishes)
	return dishes
}

// UpdateDishes updates the fill price and net change of dishes with filled orders.
func UpdateDishes() {
	dishes := AllDishes()
	for _, d := range dishes {
		d.updateFromOrders()
	}
}

// updateFromOrders sets the fill price and net change
// based on order attributes
func (d Dish) updateFromOrders() {
	if d.FillPrice != 0 {
		orders := d.Orders()
		firstOrder := orders[0]
		lastOrder := orders[len(orders)-1]

		if firstOrder.ErrorMessage == "" && lastOrder.ErrorMessage == "" {
			if firstOrder.BinanceStatus == "FILLED" || lastOrder.BinanceStatus == "FILLED" {
				var fill, net float64
				if firstOrder.BinanceStatus == "FILLED" && lastOrder.BinanceStatus == "FILLED" {
					fill = math.Abs(firstOrder.Price - lastOrder.Price)
					net = fill * firstOrder.OriginalQuantity
				} else if firstOrder.BinanceStatus == "FILLED" {
					fill = firstOrder.Price
					net = (firstOrder.Price - d.CurrentPrice) * firstOrder.OriginalQuantity
				} else {
					fill = lastOrder.Price
					net = (lastOrder.Price - d.CurrentPrice) * firstOrder.OriginalQuantity
				}
				d.SetFillPriceAndNetChange(fill, net)
			}

		}
	}

}

//SetFillPriceAndNetChange update db with fill price and net change
func (d Dish) SetFillPriceAndNetChange(fillPrice float64, netChange float64) {
	db := dbConn()
	dbPrep, err := db.Prepare("UPDATE dishes SET fill_price=?, net_change=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	dbPrep.Exec(fillPrice, netChange, d.ID)
	log.Printf("UPDATE Dish: ID: %v | Fill Price: %v | Net Change: %v", d.ID, fillPrice, netChange)

	defer db.Close()
}

// Orders returns the orders associated with a dish
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
	fmt.Printf("Dish: %v ---- Orders: %v", d, orders)
	return orders
}
