package main

import "database/sql"

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
	UPIPoints   MinMax
	UPISum      Sums
	MonthlyData []Monthlydist
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

// Amenities hold all basic amenities where most cash is transacted
type Amenities struct {
	Food    float64
	Gas     float64
	Clothes float64
	Online  float64
}

type Monthlydist struct {
	Date   string          `db:"year_month"`
	Mxdeb  float64         `db:"debsum"`
	Mxcred sql.NullFloat64 `db:"credsum"`
}
