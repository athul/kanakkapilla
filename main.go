package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// DB Schema of Sqlite is
// CREATE TABLE bank(id INTEGER PRIMARY KEY AUTOINCREMENT,tdate TEXT,
// date TEXT,desc TEXT,ref TEXT,debit FLOAT,credit FLOAT,bal FLOAT);

// Transaction struct holds a Money transaction
type Transaction struct {
	ID          int            `db:"id"`
	Tdate       string         `db:"tdate"`
	Date        string         `db:"date"`
	Description string         `db:"desc"`
	Refno       sql.NullString `db:"ref"`
	Debit       sql.NullString `db:"debit"`
	Credit      sql.NullString `db:"credit"`
	Balance     sql.NullString `db:"bal"`
}
type upiTransactions struct {
	ID          int            `db:"id"`
	Date        string         `db:"date"`
	Description string         `db:"desc"`
	Refno       sql.NullString `db:"ref"`
	Debit       sql.NullString `db:"debit"`
	Credit      sql.NullString `db:"credit"`
	Balance     sql.NullString `db:"bal"`
	Total       int
}

var (
	db  *sqlx.DB
	err error
)

func main() {
	db, err = sqlx.Open("sqlite3", "bank.db")
	if err != nil {
		log.Println(err)
	}
	trans := []Transaction{}

	if err = db.Select(&trans, `SELECT * FROM bank ORDER BY date ASC`); err != nil {
		log.Println(err)
	}
	// for _, t := range trans {
	// 	log.Println(t.Credit, "\n----")
	// 	log.Println("\n", t.Debit)
	// }
	// log.Println("----------")
	sumDebits()
	getBalance()
	getUPI()
	http.HandleFunc("/", renderTemplate)
	http.ListenAndServe(":8080", nil)
}
func renderTemplate(h http.ResponseWriter, r *http.Request) {
	temp, err := template.New("index.html").ParseFiles("index.html")
	if err != nil {
		log.Println(err)
	}
	tr := getUPI()
	if err := temp.Execute(h, tr); err != nil {
		log.Println(err)
	}

}
func getUPI() []upiTransactions {
	upiTrans := []upiTransactions{}
	if err = db.Select(&upiTrans, `SELECT id,date,desc,ref,debit,credit,bal FROM bank WHERE desc LIKE '%-UPI%' ORDER BY id DESC`); err != nil {
		log.Println("Unable to fetch UP Transactions", err)
	}
	log.Println(len(upiTrans))
	return upiTrans
}

func sumDebits() {
	var sum []float32
	if err = db.Select(&sum, `SELECT SUM(debit) FROM bank`); err != nil {
		log.Println("Unable to fetch UP Sum", err)
	}
	log.Println("Total Debited Amount", sum[0])
}

// getBalance gets the current bank balance.
// It assumes that the balance is at the last row of the table
func getBalance() {
	var balance []string
	if err = db.Select(&balance, `SELECT bal FROM bank ORDER BY id DESC LIMIT 1`); err != nil {
		log.Println("Unable to fetch UP Sum", err)
	}
	log.Println("Current Balance\t:", balance[0])
}
