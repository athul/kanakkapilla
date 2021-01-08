package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DB Schema of Sqlite is
// CREATE TABLE bank(id INTEGER PRIMARY KEY AUTOINCREMENT,tdate TEXT,
// date TEXT,desc TEXT,ref TEXT,debit FLOAT,credit FLOAT,bal FLOAT);

// Transaction struct holds a Money transaction
type Transaction struct {
	ID          int            `db:"id"`
	Tdate       string         `db:"tdate"`
	Date        string         `db:"date"`
	Description string         `db:"description"`
	Refno       sql.NullString `db:"ref"`
	Debit       sql.NullString `db:"debit"`
	Credit      sql.NullString `db:"credit"`
	Balance     float64        `db:"bal"`
}

// AllData saves all the Data from the CSV File.
type AllData struct {
	// Save all the Transcations from the account
	AllTrans []Transaction
	// Saves all the UPI transactions from the account
	UPITrans  []Transaction
	upiNos    int
	upiDebAm  float64
	CurBal    float64
	BalonDate string
	maxUpitrs int
}

var (
	db  *sqlx.DB
	err error
)

func main() {
	db, err = sqlx.Connect("postgres", "")
	if err != nil {
		log.Println(err)
	}
	trans := []Transaction{}

	if err = db.Select(&trans, `SELECT * FROM bank ORDER BY id ASC`); err != nil {
		log.Println(err)
	}

	sumDebitsfromUPI()
	upiTrans := getUPI()
	all := AllData{
		AllTrans:  trans,
		CurBal:    trans[len(trans)-1].Balance,
		BalonDate: trans[len(trans)-1].Date,
		UPITrans:  upiTrans,
	}
	http.HandleFunc("/", all.renderTemplate)
	http.ListenAndServe(":8080", nil)
}
func (a *AllData) renderTemplate(h http.ResponseWriter, r *http.Request) {
	temp, err := template.New("index.html").ParseFiles("index.html")
	if err != nil {
		log.Println(err)
	}
	// tr := getUPI()
	if err := temp.Execute(h, &a); err != nil {
		log.Println(err)
	}

}
func getUPI() []Transaction {
	upiTrans := []Transaction{}
	if err = db.Select(&upiTrans, `SELECT * FROM bank WHERE description LIKE '%-UPI%' ORDER BY id DESC`); err != nil {
		log.Println("Unable to fetch UP Transactions", err)
	}
	return upiTrans
}

func sumDebitsfromUPI() {
	var sum []float64
	if err = db.Select(&sum, `SELECT SUM(debit) FROM bank WHERE description LIKE '%-UPI%'`); err != nil {
		log.Println("Unable to fetch UP Sum", err)
	}
	log.Println("Total Debited Amount", sum[0])
}

func (t *AllData) getMinMaxupi() {

}
