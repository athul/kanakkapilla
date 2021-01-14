package main

import (
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
)

func getUPI() []Transaction {
	upiTrans := []Transaction{}
	if err = db.Select(&upiTrans, `SELECT * FROM bank WHERE description LIKE '%-UPI%' ORDER BY id DESC`); err != nil {
		log.Println("Unable to fetch UP Transactions", err)
	}
	return upiTrans
}

func (a *AllData) sumfromUPI() {
	var sum []Sums
	if err = db.Select(&sum, `SELECT SUM(debit) as debsum, SUM(credit) as credsum FROM bank WHERE description LIKE '%-UPI%'`); err != nil {
		log.Println("Unable to fetch UP Sum", err)
	}
	a.UPISum = sum[0]
	log.Println(sum)
}

func (a *AllData) getMinMaxupi() {
	minmax := []MinMax{}
	if err := db.Select(&minmax, `SELECT MAX(credit) as credmax,MAX(debit) as debmax,MIN(credit) as credmin,MIN(debit) as debmin FROM bank WHERE description LIKE '%-UPI%'`); err != nil {
		log.Println("MinMax error", err)
	}
	a.UPIPoints = minmax[0]

}
func (a *AllData) newTransaction(c echo.Context) error {
	var nildeb, nilcred interface{}
	date := c.FormValue("date")
	desc := c.FormValue("description")
	ref := c.FormValue("reference")
	debit := c.FormValue("debit")
	credit := c.FormValue("credit")
	if credit == "" {
		nilcred = nil
	} else {
		nilcred, err = strconv.ParseFloat(credit, 64)
		eros(err)
		newflt := a.CurBal + nilcred.(float64)
		log.Println("New Balance", newflt, "Old Balance", a.CurBal)
	}
	if debit == "" {
		nildeb = nil
	} else {
		nildeb, err = strconv.ParseFloat(debit, 64)
		eros(err)
		newflt := a.CurBal - nildeb.(float64)
		log.Println("New Balance", newflt, "Old Balance", a.CurBal)
	}
	log.Println(nildeb, nilcred)
	return c.String(200, "Date is "+date+"\n"+desc+"\n"+ref+"\n"+debit+"\n"+credit)
}
