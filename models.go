package main

import (
	"database/sql"
)

// DB Schema of Postgres is
// CREATE TABLE bank(id INTEGER PRIMARY KEY AUTOINCREMENT,tdate TEXT,
// date TEXT,desc TEXT,ref TEXT,debit FLOAT,credit FLOAT,bal FLOAT);

// Transaction struct holds a Money transaction
type Transaction struct {
	ID          int             `db:"id" json:"id"`
	Tdate       string          `db:"tdate" json:"tdate"`
	Date        string          `db:"date" json:"date"`
	Description string          `db:"description" json:"description"`
	Refno       sql.NullString  `db:"ref" json:"ref"`
	Debit       sql.NullFloat64 `db:"debit" json:"debit"`
	Credit      sql.NullFloat64 `db:"credit" json:"credit"`
	Balance     float64         `db:"bal" json:"bal"`
	// Tags        string          `db:"tags"`
}

// AllData saves all the Data from the CSV File.
type AllData struct {
	// Save all the Transcations from the account
	AllTrans []Transaction `json:"transactions"`
	// Saves all the UPI transactions from the account
	UPITrans []Transaction
	//CurBal is the Current Balance
	CurBal    float64
	BalonDate string
	//UPIPoints hold the Max and Min of UPI transactions,debit and credit
	MonthlyData []Monthlydist
}

//MinMax holds the Sums,Max,Mins of Debits and Credits
type MinMax struct {
	MxCredit float64 `db:"credmax" json:"credmax"`
	MxDebit  float64 `db:"debmax" json:"debmax"`
	MnCredit float64 `db:"credmin" json:"credmin"`
	MnDebit  float64 `db:"debmin" json:"debmin"`
}

// Peak holds the Maximum of UPI and Other transactions
type Peak struct {
	UPI   []MinMax `json:"upi"`
	Total []MinMax `json:"total"`
}

// Amenities hold all basic amenities where most cash is transacted
type Amenities struct {
	Food float64
	Gas  float64
	ATM  float64
	Card float64
}

// Monthlydist holds the Max of Debit and Credit from the DB
type Monthlydist struct {
	Date   string          `db:"year_month"`
	Mxdeb  sql.NullFloat64 `db:"debsum"`
	Mxcred sql.NullFloat64 `db:"credsum"`
}

// TrsVals holds the Debit and Credit value for JSON response
type TrsVals struct {
	Debit  float64 `json:"debit"`
	Credit float64 `json:"credit"`
}

// MonthTrs holds the Date and TrsVals for JSON response
type MonthTrs struct {
	Date   string  `json:"date"`
	Values TrsVals `json:"values"`
}
