package main

import (
	"encoding/json"
	"log"

	"github.com/athul/kanakkapilla/database"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

var (
	db  = database.DB
	err error
)

func eros(err error) {
	if err != nil {
		log.Println(err)
	}
}
func main() {
	db = database.InitDB()
	all := AllData{}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		data, err := json.MarshalIndent(e.Routes(), "", "  ")
		eros(err)
		return c.JSON(200, string(data))
	})
	e.POST("/ins", all.newTransaction)
	// API Group
	api := e.Group("/api")
	api.GET("/search", all.Search)
	api.GET("/month", all.RetSumsAggr)
	api.GET("/last.twenty", all.GetTrsArr)
	api.GET("/minmax", getMinMaxupi)
	api.GET("/amenities", AllAmenities)
	api.GET("/all", allTransactions)
	e.Start(":8080")
}
