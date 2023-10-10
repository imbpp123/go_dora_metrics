package graph

import (
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"

	"dora-metrics/internal/metrics/domain"
)

type LineData struct {
	Month string
	QTY   int
}

type ChartRenderOptions struct {
	Title string
	File  string
}

type EChartGenerator struct {
}

func NewEChartGenerator() *EChartGenerator {
	return &EChartGenerator{}
}

func (e *EChartGenerator) Render(data []*domain.ReleasePeriod, options ChartRenderOptions) {
	// create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme:  types.ThemeWesteros,
			Width:  "1800px",
			Height: "1200px",
		}),
		charts.WithTitleOpts(opts.Title{
			Title: options.Title,
		}))

	var xAxis []string
	yAxis := make([]opts.LineData, 0)
	for _, item := range data {
		month := item.Date().Format("2006-01")

		xAxis = append(xAxis, month)
		yAxis = append(yAxis, opts.LineData{
			Name:   month,
			Value:  item.Qty(),
			Symbol: "circle",
		})
	}

	// Put data into instance
	line.SetXAxis(xAxis).
		AddSeries("Release QTY", yAxis).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))

	f, _ := os.Create(options.File)
	_ = line.Render(f)
}
