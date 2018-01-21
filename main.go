package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"

	"github.com/pglass/sms-backup/analyze"
	"github.com/pglass/sms-backup/parse"
)

// This groups together "charts with Render function" under a type. It enables
// us to write a generic function that works with any chart type.
//
// I grepped through the go-chart source saw no equivalent interface.
type RenderableChart interface {
	Render(rp chart.RendererProvider, w io.Writer) error
}

const (
	CHART_HEIGHT = 750
	CHART_WIDTH  = 960
	CHART_DPI    = 110.0
)

var (
	FILENAME  string
	OUTFILE   string
	MY_NUMBER string

	CHART_TYPE string

	X_AXIS chart.XAxis
	Y_AXIS chart.YAxis
)

func init() {
	flag.StringVar(&FILENAME, "f", "", "The XML file containing your SMS backups")
	flag.StringVar(&OUTFILE, "o", "out.png", "The output image")
	flag.StringVar(&MY_NUMBER, "n", "", "My phone number. Used to determine if MMS messages are incoming")
	flag.StringVar(&CHART_TYPE, "t", "",
		"One of: messagesPerDay, messagesPerWeek",
	)

	X_AXIS = chart.XAxis{
		Style: chart.Style{
			Show: true,
		},
	}
	Y_AXIS = chart.YAxis{
		Style: chart.Style{
			Show: true,
		},
	}

}

func main() {
	flag.Parse()
	log.SetFlags(0)

	if FILENAME != "" {
		if doc, err := parse.ParseXML(FILENAME); err != nil {
			log.Fatal(err)
		} else {
			for i := range doc.MMSes {
				doc.MMSes[i].SetMyNumber(MY_NUMBER)
			}
			MakeChart(doc, CHART_TYPE)
		}
	} else {
		log.Fatal("Need a filename")
	}
}

func MakeChart(doc parse.Document, chartType string) {
	analyzer := analyze.MakeAnalyzer(doc)
	analyzer.Run()

	var plot RenderableChart
	switch chartType {
	case "messagesPerWeek":
		plot = GetTimeSeriesChart(
			[]chart.TimeSeries{
				// MakeTimeSeries(analyzer.MessagesPerWeek, "Messages per week"),
				MakeTimeSeries(analyzer.IncomingMessagesPerWeek, "Incoming Messages per week"),
				MakeTimeSeries(analyzer.OutgoingMessagesPerWeek, "Outgoing Messages per week"),
			},
		)
	case "messagesPerDay":
		plot = GetTimeSeriesChart(
			[]chart.TimeSeries{
				// MakeTimeSeries(analyzer.MessagesPerDay, "Messages per day"),
				MakeTimeSeries(analyzer.IncomingMessagesPerDay, "Incoming Messages per day"),
				MakeTimeSeries(analyzer.OutgoingMessagesPerDay, "Outgoing Messages per day"),
			},
		)
	case "incomingMessageLengths":
		plot = GetHistogramChart(analyzer.IncomingMessageLengths, "Incoming Message Lengths")
	case "outgoingMessageLengths":
		plot = GetHistogramChart(analyzer.OutgoingMessageLengths, "Outgoing Message Lengths")
	case "messagesTimeOfDay":
		c := GetTimeSeriesChart(
			[]chart.TimeSeries{
				MakeTimeSeriesScatter(analyzer.MessagesTimeOfDay, "Messages plotted by hour of day"),
			},
		)
		c.YAxis = chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
			Ticks: []chart.Tick{
				{0.0, "12:00am"},
				{2.0, "2:00am"},
				{4.0, "4:00am"},
				{6.0, "6:00am"},
				{8.0, "8:00am"},
				{10.0, "10:00am"},
				{12.0, "12:00pm"},
				{14.0, "2:00pm"},
				{16.0, "4:00pm"},
				{18.0, "6:00pm"},
				{20.0, "8:00pm"},
				{22.0, "10:00pm"},
				{24.0, "12:00am"},
			},
		}
		plot = c
	default:
		log.Fatalf("Unsupported chart type: %v", chartType)
	}

	RenderPlot(plot, OUTFILE)
}

func MakeTimeSeries(m map[time.Time]float64, name string) chart.TimeSeries {
	keys, values := analyze.SplitMapSorted(m)
	return chart.TimeSeries{Name: name, XValues: keys, YValues: values}
}

func MakeTimeSeriesScatter(m map[time.Time]float64, name string) chart.TimeSeries {
	viridisByY := func(xr, yr chart.Range, index int, x, y float64) drawing.Color {
		return chart.Viridis(y, yr.GetMin(), yr.GetMax())
	}

	keys, values := analyze.SplitMapSorted(m)
	return chart.TimeSeries{
		Name:    name,
		XValues: keys,
		YValues: values,
		Style: chart.Style{
			Show:             true,
			StrokeWidth:      chart.Disabled,
			DotWidth:         5,
			DotColorProvider: viridisByY,
		},
	}
}

func GetTimeSeriesChart(time_series []chart.TimeSeries) chart.Chart {
	series := make([]chart.Series, len(time_series))
	for i, val := range time_series {
		series[i] = val
	}
	result := chart.Chart{
		XAxis:  X_AXIS,
		YAxis:  Y_AXIS,
		Series: series,

		Height: CHART_HEIGHT,
		Width:  CHART_WIDTH,
		DPI:    CHART_DPI,
	}
	result.Elements = []chart.Renderable{chart.Legend(&result)}
	return result
}

func GetHistogramChart(histogram analyze.Histogram, title string) chart.BarChart {
	data := histogram.GetMap()

	keys, vals := analyze.SplitMapSortedInt64(data)

	chart_vals := []chart.Value{}
	for i, k := range keys {
		chart_vals = append(chart_vals, chart.Value{Value: vals[i], Label: fmt.Sprintf("%v", k)})
	}
	result := chart.BarChart{
		Title:      title,
		TitleStyle: chart.StyleShow(),
		XAxis:      chart.Style{Show: true},
		YAxis:      Y_AXIS,
		Bars:       chart_vals,

		Height: CHART_HEIGHT,
		Width:  CHART_WIDTH,
		DPI:    CHART_DPI,
	}
	// result.Elements = []chart.Renderable{chart.Legend(&result)}
	return result
}

func GetBarChart(keys []time.Time, values []float64) chart.BarChart {
	chart_vals := []chart.Value{}
	for i, _ := range keys {
		chart_vals = append(chart_vals, chart.Value{Value: float64(values[i]), Label: ""})
	}
	result := chart.BarChart{
		XAxis: chart.Style{Show: true},
		YAxis: Y_AXIS,
		Bars:  chart_vals,

		Height: CHART_HEIGHT,
		Width:  CHART_WIDTH,
		DPI:    CHART_DPI,
	}
	// result.Elements = []chart.Renderable{chart.Legend(&result)}
	return result
}

func RenderPlot(plot RenderableChart, filename string) {
	buf := bytes.NewBuffer([]byte{})

	mode := InferRendererProviderFromFilename(filename)
	if err := plot.Render(mode, buf); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Wrote file %v", filename)
	}
}

func InferRendererProviderFromFilename(filename string) chart.RendererProvider {
	switch filename[len(filename)-4:] {
	case ".svg":
		return chart.SVG
	case ".png":
		return chart.PNG
	default:
		log.Fatalf("Unsupported file extension: %v (try '.png' or '.svg'", filename)
	}
	return nil
}
