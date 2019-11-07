package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Recipe struct {
	ID        int
	Name      string
	Symbol    string
	Side      string
	GainRatio float64
	LossRatio float64
	Amount    float64
	Frequency int
}

func (r Recipe) CreateOrders() string {
	currentPrice := GetPrice("BNBBTC")

	highPrice := MultiplyPrice(currentPrice, r.GainRatio)
	lowPrice := MultiplyPrice(currentPrice, r.LossRatio)
	stopPrice := MultiplyPrice(lowPrice, 0.99)
	CreateOrder(highPrice)
	CreateStopLossLimitOrder(lowPrice, stopPrice)
	return "ok"
}

type Dish struct {
	ID           int
	RecipeID     int
	PairOne      string
	PairTwo      string
	CurrentPrice float64
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "trade"
	dbPass := "tradebot"
	dbname := "trade_bot"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbname)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM recipes ORDER BY id ASC")
	if err != nil {
		panic(err.Error())
	}
	recipe := Recipe{}
	recipes := []Recipe{}
	for selDB.Next() {
		var id int
		var name, symbol, side string
		var gainRatio, lossRatio, amount float64
		var frequency int
		err = selDB.Scan(&id, &name, &symbol, &side, &gainRatio, &lossRatio, &amount, &frequency)
		if err != nil {
			panic(err.Error())
		}
		recipe.ID = id
		recipe.Name = name
		recipe.Symbol = symbol
		recipe.Side = side
		recipe.GainRatio = gainRatio
		recipe.LossRatio = lossRatio
		recipe.Amount = amount
		recipe.Frequency = frequency
		recipes = append(recipes, recipe)
	}
	fmt.Println("Index Recipes:", recipes)
	tmpl.ExecuteTemplate(w, "Index", recipes)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM recipes WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	recipe := Recipe{}
	for selDB.Next() {
		var id int
		var name, symbol, side string
		var gainRatio, lossRatio, amount float64
		var frequency int
		err = selDB.Scan(&id, &name, &symbol, &side, &gainRatio, &lossRatio, &amount, &frequency)
		if err != nil {
			panic(err.Error())
		}
		recipe.ID = id
		recipe.Name = name
		recipe.Symbol = symbol
		recipe.Side = side
		recipe.GainRatio = gainRatio
		recipe.LossRatio = lossRatio
		recipe.Amount = amount
		recipe.Frequency = frequency
	}
	tmpl.ExecuteTemplate(w, "Show", recipe)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM recipes WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	recipe := Recipe{}
	for selDB.Next() {
		var id int
		var name, symbol, side string
		var gainRatio, lossRatio, amount float64
		var frequency int
		err = selDB.Scan(&id, &name, &symbol, &side, &gainRatio, &lossRatio, &amount, &frequency)
		if err != nil {
			panic(err.Error())
		}
		recipe.ID = id
		recipe.Name = name
		recipe.Symbol = symbol
		recipe.Side = side
		recipe.GainRatio = gainRatio
		recipe.LossRatio = lossRatio
		recipe.Amount = amount
		recipe.Frequency = frequency
	}
	tmpl.ExecuteTemplate(w, "Edit", recipe)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("Name")
		symbol := r.FormValue("Symbol")
		side := r.FormValue("Side")
		gainRatio := r.FormValue("GainRatio")
		lossRatio := r.FormValue("LossRatio")
		amount := r.FormValue("Amount")
		frequency := r.FormValue("Frequency")
		insForm, err := db.Prepare("INSERT INTO recipes(name, symbol, side, gain_ratio, loss_ratio, amount, frequency) VALUES(?,?,?,?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, symbol, side, gainRatio, lossRatio, amount, frequency)
		log.Println("INSERT: name: " + name + " | Symbol: " + symbol + " | Side: " + side)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		PairOne := r.FormValue("PairOne")
		PairTwo := r.FormValue("PairTwo")
		GainRatio := r.FormValue("GainRatio")
		LossRatio := r.FormValue("LossRatio")
		Amount := r.FormValue("Amount")
		Frequency := r.FormValue("Frequency")
		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE recipes SET name=?, PairOne=?, PairTwo=?, GainRatio=?, LossRatio=?, Amount=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, PairOne, PairTwo, GainRatio, LossRatio, Amount, Frequency, id)
		log.Println("UPDATE: name: " + name + " | PairOne: " + PairOne + " | PairTwo: " + PairTwo)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	recipe := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM recipes WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(recipe)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {

	GetPrice("BNBBTC")
	//
	// db := dbConn()
	// selDB, err := db.Query("SELECT * FROM recipes ORDER BY id DESC LIMIT 1")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// recipe := Recipe{}
	// for selDB.Next() {
	// 	var id int
	// 	var name, PairOne, PairTwo string
	// 	var GainRatio, LossRatio, Amount float64
	// 	var Frequency int
	// 	err = selDB.Scan(&id, &name, &PairOne, &PairTwo, &GainRatio, &LossRatio, &Amount, &Frequency)
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	// 	recipe.Id = id
	// 	recipe.Name = name
	// 	recipe.PairOne = PairOne
	// 	recipe.PairTwo = PairTwo
	// 	recipe.GainRatio = GainRatio
	// 	recipe.LossRatio = LossRatio
	// 	recipe.Amount = Amount
	// 	recipe.Frequency = Frequency
	// }

	// defer db.Close()
	//
	// fmt.Println("Recipe: %v", recipe)
	// recipe.CreateOrders()

	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)

}
