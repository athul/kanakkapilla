package database

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //postgres driver
)

var (
	//DB is the sqlx db type
	DB  *sqlx.DB
	err error
)

// CSV Data will be of the form
// +----+------------+------------+-----------------------------+-----------------------+--------+--------+---------+
// | id | tdate      | date       | desc                        | ref                   | debit  | credit | bal     |
// +----+------------+------------+-----------------------------+-----------------------+--------+--------+---------+
// | 1  | 1 Feb 2020 | 1 Feb 2020 | TO TRANSFER TO SOMEONE      | TRANSFER TO 123456789 | 100.00 |        | 6237.15 |
// +----+------------+------------+-----------------------------+-----------------------+--------+--------+---------+
// | 2  | 8 Feb 2020 | 8 Feb 2020 | TO TRANSFER TO SOMEONE ELSE | TRANSFER TO 123456789 |        | 100.00 | 6337.15 |
// +----+------------+------------+-----------------------------+-----------------------+--------+--------+---------+

//InitDB initializes the DB
func InitDB() *sqlx.DB {
	// DB Schema
	//CREATE TABLE IF NOT EXISTS bank(
	// id SERIAL PRIMARY KEY,
	// tdate DATE,
	// date DATE,
	// description TEXT,
	// ref TEXT DEFAULT NULL,
	// debit NUMERIC DEFAULT NULL,
	// credit NUMERIC DEFAULT NULL,
	// bal NUMERIC
	// );`
	DB, err = sqlx.Connect("postgres", "user=athul password=splendor sslmode=disable") //os.Getenv("pgurl"))
	if err != nil {
		log.Println(err)
	}
	if err = DB.Ping(); err != nil {
		log.Println("Ping Error", err)
	}
	return DB
}

// readjson reads the transactions in JSON and parses it
func readjson(file string) string {
	file_data, err := os.ReadFile(file)
	if err != nil {
		log.Println("JSON Reading Error", err)
	}
	return string(file_data)
}

// InserttoDB reads the JSON file and inserts the data
// to Postgres via `json_populate_recordset` function of postgres
func InserttoDB(file string) {
	db := InitDB() //Initialize DB for Inserting Data
	json_file := readjson(file)
	ins := `INSERT INTO bank SELECT * FROM json_populate_recordset(NULL::bankjs,$1);`
	ls, err := db.Exec(ins, json_file)
	if err != nil {
		log.Println("Error Inserting to table", err)
	}
	log.Println(ls.RowsAffected())
}
