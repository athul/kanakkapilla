package main

import (
	"log"

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
func newTransaction(c echo.Context) error {
	ll := c.FormValue("date")
	log.Println(ll)
	return c.String(200, "Date is "+ll)
}
