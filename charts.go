package main

import (
	"bytes"
	"fmt"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/labstack/echo/v4"
)

func (a *AllData) genAllGraphOptions(deb, bal bool) []opts.BarData {
	items := make([]opts.BarData, 0)
	if deb {
		for _, dt := range a.AllTrans {
			items = append(items, opts.BarData{Value: dt.Debit.Float64})
		}
	} else if bal {
		for _, dt := range a.AllTrans {
			items = append(items, opts.BarData{Value: dt.Balance})
		}
	} else {
		for _, dt := range a.AllTrans {
			items = append(items, opts.BarData{Value: dt.Credit.Float64})
		}
	}
	return items
}
func (a *AllData) genChart(c echo.Context) error {
	var b bytes.Buffer
	line := charts.NewBar()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ChartThemeRiver, Width: "1300px"}),
		charts.WithTitleOpts(opts.Title{
			Title:    fmt.Sprintf("Till %s", a.BalonDate),
			Subtitle: "All Transactions",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: true,
		}))
	var dates []string
	for _, dt := range a.AllTrans {
		dates = append(dates, dt.Date)
	}
	// Put data into instance
	line.SetXAxis(dates).
		AddSeries("Debit", a.genAllGraphOptions(true, false)).
		AddSeries("Credit", a.genAllGraphOptions(false, false)).
		AddSeries("Balance", a.genAllGraphOptions(false, true)).
		SetSeriesOptions(charts.WithMarkPointNameTypeItemOpts(
			opts.MarkPointNameTypeItem{Name: "Maximum", Type: "max"},
			opts.MarkPointNameTypeItem{Name: "Average", Type: "average"},
			opts.MarkPointNameTypeItem{Name: "Minimum", Type: "min"},
		), charts.WithMarkPointNameCoordItemOpts(opts.MarkPointNameCoordItem{
			Name: "data",
		}))
	line.Render(&b)
	return c.HTML(200, b.String())
}
