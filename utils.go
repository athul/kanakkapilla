package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strings"

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
func (a *AllData) renderSearch(e echo.Context) error {
	descName := e.FormValue("query")
	trs := []Transaction{}
	var (
		b bytes.Buffer
	)
	query := fmt.Sprintf(`SELECT * FROM bank WHERE description LIKE %s OR description LIKE %s ORDER BY id DESC`, "'%"+strings.ToUpper(descName)+"%'", "'%"+descName+"%'")
	log.Println("Query", query)
	err := db.Select(&trs, query)
	if err != nil {
		log.Println("Get Error", err)
	}
	log.Println(trs)
	// for rows.Next() {
	// 	err := rows.Scan(&trs)
	// 	eros(err)

	// }
	temp, err := template.New("table").Funcs(funcMap).Parse(`
	<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Find Transaction</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@exampledev/new.css@1.1.2/new.min.css">
</head>
<body>
    <table>
        <tr>
            <th>ID</th>
            <th>TDate</th>
            <th>Date</th>
            <th>Description</th>
            <th>Ref</th>
            <th>Debit</th>
            <th>Credit</th>
            <th>Balance</th>
        </tr>
        {{ range .}}
            <tr>
                <td>{{ .ID }}</td>
                <td>{{ .Tdate }}</td>
                <td>{{ .Date }}</td>
                <td>{{ .Description }}</td>
                <td>{{ .Refno.String }}</td>
                <td>{{ .Debit.Float64| humanize | currit }}</td>
                <td>{{ .Credit.Float64| humanize | currit }}</td>
                <td>{{ .Balance| humanize | currit }}</td>
            </tr>
        {{ end}}
    </table>
</body>
<style>
body{
	max-width:max-content !important;
	}
</style>
</html>`)
	eros(err)
	temp.Execute(&b, trs)
	return e.HTML(200, b.String())
}
