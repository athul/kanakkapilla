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
	trans := []Transaction{}

	if err = db.Select(&trans, `SELECT * FROM bank ORDER BY id ASC`); err != nil {
		log.Println(err)
	}

	all := AllData{
		AllTrans:  trans,
		CurBal:    trans[len(trans)-1].Balance,
		BalonDate: trans[len(trans)-1].Date,
		UPITrans:  getUPI(),
	}
	all.getMinMaxupi()
	all.sumfromUPI()
	e := echo.New()
	e.GET("/", all.renderIndexTemplate)
	e.GET("/all", all.renderTableTemplate)
	e.GET("/all.graph", all.genChart)
	e.POST("/search", all.renderSearch)
	e.Start(":8080")
}
