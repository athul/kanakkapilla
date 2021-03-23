package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

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
	var (
		nildeb, nilcred interface{}
		newflt          float64
	)
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
		newflt = a.CurBal + nilcred.(float64)
		log.Println("New Balance", newflt, "Old Balance", a.CurBal)
	}
	if debit == "" {
		nildeb = nil
	} else {
		nildeb, err = strconv.ParseFloat(debit, 64)
		eros(err)
		newflt = a.CurBal - nildeb.(float64)
		log.Println("New Balance", newflt, "Old Balance", a.CurBal)
	}
	insStmt := `INSERT INTO bank VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	res := db.MustExec(insStmt, len(a.AllTrans)+1, date, date, desc, ref, nildeb, nilcred, fmt.Sprintf("%.2f", newflt))
	log.Println(res.RowsAffected())
	return c.String(200, "Date is "+date+"\n"+desc+"\n"+ref+"\n"+debit+"\n"+credit+"\n"+fmt.Sprintf("%.2f", newflt))
}

//GetTrsArr returns a slice of the debit/credit transactions
func (a *AllData) GetTrsArr(mode string) []float64 {
	trs := make([]float64, 0)
	switch mode {
	case "deb":
		for _, dt := range a.AllTrans {
			trs = append(trs, dt.Debit.Float64)
		}
	case "cred":
		for _, dt := range a.AllTrans {
			trs = append(trs, dt.Credit.Float64)
		}
	}
	return trs[:10]
}

func (a *AllData) GetDates() []string {
	var dates []string
	for _, dt := range a.AllTrans {
		t, err := time.Parse(time.RFC3339, dt.Date)
		eros(err)
		dates = append(dates, t.Format("Jan 2, 2006"))
	}

	return dates[:10]
}
