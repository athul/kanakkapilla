package csv2pg

import (
	"log"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //postgres driver
)

var (
	//DB is the sqlx db type
	DB  *sqlx.DB
	err error
)

// Transaction holds the structure of the csv and db
type Transaction struct {
	ID          int    `db:"id" csv:"id"`
	Tdate       string `db:"tdate" csv:"tdate"`
	Date        string `db:"date" csv:"date"`
	Description string `db:"description" csv:"desc"`
	Refno       string `db:"ref" csv:"ref"`
	Debit       string `db:"debit" csv:"debit"`
	Credit      string `db:"credit" csv:"credit"`
	Balance     string `db:"bal" csv:"bal"`
}

//InitDB initializes the DB
func InitDB() *sqlx.DB {
	// var tableSchema = `DROP TABLE bank;CREATE TABLE IF NOT EXISTS bank(
	// id SERIAL PRIMARY KEY,
	// tdate TEXT,
	// date TEXT,
	// description TEXT,
	// ref TEXT DEFAULT NULL,
	// debit NUMERIC DEFAULT NULL,
	// credit NUMERIC DEFAULT NULL,
	// bal NUMERIC
	// );`
	DB, err = sqlx.Connect("postgres", "")
	if err != nil {
		log.Println(err)
	}
	if err = DB.Ping(); err != nil {
		log.Println("Ping Error", err)
	}
	// result, err := db.Exec(tableSchema)
	// if err != nil {
	// 	log.Println("Table Creation Error", err)
	// }
	// log.Println(result.RowsAffected())
	// transactions := handleCSV()
	// csvtopostgres(transactions)
	return DB
}

func handleCSV() []Transaction {
	csvFile, err := os.Open("bank-str.csv")
	if err != nil {
		log.Println(err)
	}
	defer csvFile.Close()
	trans := []Transaction{}
	if err := gocsv.UnmarshalFile(csvFile, &trans); err != nil {
		log.Println("CSV Unmarshal Error", err)
	}
	return trans
}

func csvtopostgres(trs []Transaction) {
	var cred, deb interface{}
	insQuery := `INSERT INTO bank VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`
	for _, t := range trs {
		if len(t.Credit) == 0 {
			cred = nil
		} else {
			cred = t.Credit
		}
		if len(t.Debit) == 0 {
			deb = nil
		} else {
			deb = t.Debit
		}
		ls, err := DB.Exec(insQuery, t.ID, t.Tdate, t.Date, t.Description, t.Refno, deb, cred, t.Balance)
		if err != nil {
			log.Println("Execution error", err)
		}
		log.Println(ls.RowsAffected())
	}
}
