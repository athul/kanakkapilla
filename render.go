package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/labstack/echo/v4"
)

var funcMap = template.FuncMap{
	"humanize": func(fl float64) string {
		return humanize.Commaf(fl)
	},
	"currit": func(as string) string {
		return fmt.Sprintf("₹ %s", as)
	},
	"hdate": func(d string) string {
		t, err := time.Parse(time.RFC3339, d)
		eros(err)
		return t.Format("Jan 2, 2006")
	},
}

func (a *AllData) renderIndexTemplate(c echo.Context) error {
	var b bytes.Buffer

	temp, err := template.New("index.html").Funcs(funcMap).ParseGlob("templates/index.html")
	if err != nil {
		log.Println(err)
	}
	// tr := getUPI()
	if err := temp.Execute(&b, &a); err != nil {
		log.Println(err)
	}
	return c.HTML(200, b.String())
}
func (a *AllData) renderUPITemplate(c echo.Context) error {
	var b bytes.Buffer

	temp, err := template.New("upi.html").Funcs(funcMap).ParseGlob("templates/upi.html")
	if err != nil {
		log.Println(err)
	}
	// tr := getUPI()
	if err := temp.Execute(&b, &a); err != nil {
		log.Println(err)
	}
	return c.HTML(200, b.String())
}
func (a *AllData) renderTableTemplate(c echo.Context) error {
	var b bytes.Buffer

	temp, err := template.New("table.html").Funcs(funcMap).ParseGlob("templates/table.html")
	if err != nil {
		log.Println(err)
	}
	// tr := getUPI()
	if err := temp.Execute(&b, &a); err != nil {
		log.Println(err)
	}
	return c.HTML(200, b.String())
}
func (a *AllData) renderSearch(e echo.Context) error {
	descName := e.FormValue("query")
	trs := []Transaction{}
	var b bytes.Buffer
	query := fmt.Sprintf(`SELECT * FROM bank WHERE description LIKE %s OR description LIKE %s ORDER BY date DESC`, "'%"+strings.ToUpper(descName)+"%'", "'%"+descName+"%'")
	log.Println("Query", query)
	err := db.Select(&trs, query)
	if err != nil {
		log.Println("Get Error", err)
	}
	log.Println(trs)
	temp, err := template.New("search-table.html").Funcs(funcMap).ParseGlob("templates/search-table.html")
	eros(err)
	if err := temp.Execute(&b, trs); err != nil {
		log.Println(err)
	}
	return e.HTML(200, b.String())
}
