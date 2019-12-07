package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var tmpl = template.Must(template.ParseGlob("views/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM recipes ORDER BY id ASC")
	if err != nil {
		panic(err.Error())
	}
	recipe := Recipe{}
	recipes := []Recipe{}
	for selDB.Next() {
		var (
			id                             int
			name, symbol, side             string
			gainRatio, lossRatio, quantity float64
			frequency                      int
			active                         bool
		)

		err = selDB.Scan(&id, &name, &symbol, &side, &gainRatio, &lossRatio, &quantity, &frequency, &active)
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
		var (
			id                             int
			name, symbol, side             string
			gainRatio, lossRatio, quantity float64
			frequency                      int
			active                         bool
		)
		err = selDB.Scan(&id, &name, &symbol, &side, &gainRatio, &lossRatio, &quantity, &frequency, &active)
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
		var (
			id                             int
			name, symbol, side             string
			gainRatio, lossRatio, quantity float64
			frequency                      int
			active                         bool
		)
		err = selDB.Scan(&id, &name, &symbol, &side, &gainRatio, &lossRatio, &quantity, &frequency, &active)
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
		quantity := r.FormValue("Quantity")
		frequency := r.FormValue("Frequency")
		active := r.FormValue("Active")
		insForm, err := db.Prepare("INSERT INTO recipes(name, symbol, side, gain_ratio, loss_ratio, quantity, frequency, active) VALUES(?,?,?,?,?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, symbol, side, gainRatio, lossRatio, quantity, frequency, active)
		log.Println("INSERT: name: " + name + " | Symbol: " + symbol + " | Side: " + side)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("Name")
		symbol := r.FormValue("Symbol")
		side := r.FormValue("Side")
		gainRatio := r.FormValue("GainRatio")
		lossRatio := r.FormValue("LossRatio")
		quantity := r.FormValue("Quantity")
		frequency := r.FormValue("Frequency")
		active := r.FormValue("Active")
		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE recipes SET name=?, symbol=?, side=?, gain_ratio=?, loss_ratio=?, quantity=?, frequency=?, active=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, symbol, side, gainRatio, lossRatio, quantity, frequency, active, id)
		log.Println("UPDATE: name: " + name + " | Symbol: " + symbol + " | Side: " + side)
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

func Cook(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM recipes WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	recipe := Recipe{}
	for selDB.Next() {
		var (
			id                             int
			name, symbol, side             string
			gainRatio, lossRatio, quantity float64
			frequency                      int
			active                         bool
		)
		err = selDB.Scan(&id, &name, &symbol, &side, &gainRatio, &lossRatio, &quantity, &frequency, &active)
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
	}
	recipe.CookDish()
	tmpl.ExecuteTemplate(w, "Show", recipe)
	defer db.Close()
}

func serveWeb() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.HandleFunc("/cook", Cook)
	http.ListenAndServe(":8080", nil)
}
