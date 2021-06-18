package main

import (
	"log"
	"net/http"

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
	e.GET("/", all.renderIndexTemplate)
	e.GET("/all", all.renderTableTemplate)
	e.GET("/upi", all.renderUPITemplate)
	e.File("/new", "templates/insert.html")
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

func allTransactions(c echo.Context) error {
	toDate := c.QueryParam("toDate")
	fromDate := c.QueryParam("fromDate")
	trans := []Transaction{}
	if toDate != "" && fromDate != "" {
		err = db.Select(&trans, `SELECT * FROM bank WHERE date BETWEEN $1 AND $2 ORDER BY id DESC`, fromDate, toDate)
		eros(err)
	} else {
		err = db.Select(&trans, `SELECT * FROM bank ORDER BY id DESC`)
		eros(err)
	}
	return c.JSON(http.StatusOK, trans)
}
