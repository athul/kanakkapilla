package main

import (
	"log"

	"github.com/athul/kanakkapilla/csv2pg"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

var (
	db  = csv2pg.DB
	err error
)

func eros(err error) {
	if err != nil {
		log.Println(err)
	}
}
func main() {
	db = csv2pg.InitDB()
	all := AllData{}
	// all.fetchAlltrs()
	// all.getMinMaxupi()
	// all.sumfromUPI()
	// all.CalcMonthlyMax()
	e := echo.New()
	e.GET("/", all.renderIndexTemplate)
	e.GET("/all", all.renderTableTemplate)
	// e.GET("/all.graph", all.genChart)
	e.GET("/upi", all.renderUPITemplate)
	e.File("/new", "templates/insert.html")
	e.POST("/search", all.renderSearch)
	e.POST("/ins", all.newTransaction)
	e.Start(":8080")
}

func (a *AllData) fetchAlltrs() {
	trans := []Transaction{}
	err = db.Select(&trans, `SELECT * FROM bank ORDER BY id DESC`)
	eros(err)
	a.AllTrans = trans
	a.CurBal = trans[0].Balance
	a.BalonDate = trans[0].Date
	a.UPITrans = getUPI()
}
