package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

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
