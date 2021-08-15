package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func getMinMaxupi(c echo.Context) error {
	peak := Peak{}
	if err := db.Select(&peak.UPI, `SELECT MAX(credit) as credmax,MAX(debit) as debmax,MIN(credit) as credmin,MIN(debit) as debmin FROM bank WHERE description LIKE '%-UPI%'`); err != nil {
		log.Println("MinMax error", err)
	}
	if err := db.Select(&peak.Total, `SELECT MAX(credit) as credmax,MAX(debit) as debmax,MIN(credit) as credmin,MIN(debit) as debmin FROM bank`); err != nil {
		log.Println("MinMax error", err)
	}
	return c.JSON(http.StatusOK, peak)
}
func (a *AllData) getAlltrs() {
	trs := []Transaction{}
	err = db.Select(&trs, `SELECT * FROM bank;`)
	eros(err)
	a.AllTrans = trs
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
func (a *AllData) GetTrsArr(c echo.Context) error {
	lasttrs := []MonthTrs{}
	a.getAlltrs()
	for _, dt := range a.AllTrans {
		t, err := time.Parse(time.RFC3339, dt.Date)
		eros(err)
		date := t.Format("Jan 2, 2006")
		lasttrs = append(lasttrs, MonthTrs{
			Date: date,
			Values: TrsVals{
				Debit:  dt.Debit.Float64,
				Credit: dt.Credit.Float64,
			},
		})
	}
	log.Println(lasttrs)
	return c.JSON(200, lasttrs[len(lasttrs)-20:])
}

// CalcMonthlyMax calculates the Sum of Debit and Credit per month
func (a *AllData) CalcMonthlyMax() {
	var mdist []Monthlydist
	if err := db.Select(&mdist, `SELECT to_char(date,'YY-MM(MONTH)') AS year_month, SUM(debit) AS debsum, sum(credit) AS credsum FROM bank GROUP BY year_month ORDER BY year_month;`); err != nil {
		log.Println("Monthly Data Aggregation Error:", err)
	}
	log.Println("Monthly Data \n", mdist)
	a.MonthlyData = mdist
}

// RetSumsAggr returns Monthly Total Transactions as JSON
func (a *AllData) RetSumsAggr(c echo.Context) error {
	trs := []MonthTrs{}
	a.CalcMonthlyMax()
	for _, dt := range a.MonthlyData {
		trs = append(trs, MonthTrs{
			Date: strings.ReplaceAll(dt.Date, " ", ""),
			Values: TrsVals{
				Debit:  dt.Mxdeb.Float64,
				Credit: dt.Mxcred.Float64,
			},
		})
	}
	return c.JSON(200, trs)
}
func (a *AllData) Search(e echo.Context) error {
	descName := e.QueryParam("query")
	trs := []Transaction{}
	query := fmt.Sprintf(`SELECT * FROM bank WHERE description LIKE %s OR description LIKE %s ORDER BY date DESC`, "'%"+strings.ToUpper(descName)+"%'", "'%"+descName+"%'")
	log.Println("Query", query)
	err := db.Select(&trs, query)
	if err != nil {
		log.Println("Get Error", err)
	}
	return e.JSON(200, trs)
}
func (am *Amenities) Calcfuel() {
	var fuel []float64
	if err := db.Select(&fuel, `SELECT SUM(debit) FROM bank WHERE description LIKE '%FUEL%' OR description LIKE '%OYOOS%';`); err != nil {
		log.Println(err)
	}
	am.Gas = fuel[0]
}
func (am *Amenities) Calcfood() {
	var food []float64
	if err := db.Select(&food, `SELECT sum(debit) FROM bank WHERE description LIKE '%ROYAL%' OR description LIKE '%RETAIL%' OR description LIKE '%swiggy%' OR description LIKE '%PEEDIKA%';`); err != nil {
		log.Println(err)
	}
	am.Food = food[0]
}
func AllAmenities(c echo.Context) error {
	a := Amenities{}
	a.Calcfood()
	a.Calcfuel()
	log.Println(a)
	return c.JSON(200, a)
}

func allTransactions(c echo.Context) error {
	toDate := c.QueryParam("toDate")
	fromDate := c.QueryParam("fromDate")
	trans := []Transaction{}
	if toDate != "" && fromDate != "" {
		err = db.Select(&trans, `SELECT * FROM bank WHERE date BETWEEN $1 AND $2`, fromDate, toDate)
		eros(err)
	} else {
		err = db.Select(&trans, `SELECT * FROM bank`)
		eros(err)
	}
	return c.JSON(http.StatusOK, trans)
}
