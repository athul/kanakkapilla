package main

import (
	"database/sql"
	"log"

	"github.com/athul/kanakkapilla/csv2pg"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

// DB Schema of Postgres is
// CREATE TABLE bank(id INTEGER PRIMARY KEY AUTOINCREMENT,tdate TEXT,
// date TEXT,desc TEXT,ref TEXT,debit FLOAT,credit FLOAT,bal FLOAT);

// Transaction struct holds a Money transaction
type Transaction struct {
	ID          int             `db:"id"`
	Tdate       string          `db:"tdate"`
	Date        string          `db:"date"`
	Description string          `db:"description"`
	Refno       sql.NullString  `db:"ref"`
	Debit       sql.NullFloat64 `db:"debit"`
	Credit      sql.NullFloat64 `db:"credit"`
	Balance     float64         `db:"bal"`
}

// AllData saves all the Data from the CSV File.
type AllData struct {
	// Save all the Transcations from the account
	AllTrans []Transaction
	// Saves all the UPI transactions from the account
	UPITrans []Transaction
	upiNos   int
	upiDebAm float64
	//CurBal is the Current Balance
	CurBal    float64
	BalonDate string
	//UPIPoints hold the Max and Min of UPI transactions,debit and credit
	UPIPoints MinMax
	UPISum    Sums
}

//MinMax holds the Max and Mins of Debits and Credits
type MinMax struct {
	MxCredit float64 `db:"credmax"`
	MxDebit  float64 `db:"debmax"`
	MnCredit float64 `db:"credmin"`
	MnDebit  float64 `db:"debmin"`
}

// Sums holds the sum of Debited and Credit Amounts
type Sums struct {
	DebSum  float64 `db:"debsum"`
	CredSum float64 `db:"credsum"`
}

var (
	db  = csv2pg.DB
	err error
)

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
	e.Start(":8080")
}

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
