package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"strconv"

	chart "github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

func plot(lines ...[][]float64) {
	fmt.Printf("Plotting Graph\n")

	var charts []chart.Series
	for i, line := range lines {
		currentLine := chart.Series(
			chart.ContinuousSeries{
				Name: "Time Series",
				Style: chart.Style{
					Show: true,
				},
				XValues: line[0],
				YValues: line[1],
			},
		)
		currentAnotation := chart.AnnotationSeries{
			Annotations: []chart.Value2{
				{XValue: line[0][len(line[0])-1], YValue: line[1][len(line[1])-1], Label: strconv.Itoa(i)},
			},
		}

		charts = append(charts, currentLine)
		charts = append(charts, currentAnotation)

	}

	graph := chart.Chart{
		Background: chart.Style{
			Padding: chart.Box{
				Top:    50,
				Left:   25,
				Right:  25,
				Bottom: 10,
			},
			FillColor: drawing.ColorFromHex("efefef"),
		},
		XAxis: chart.XAxis{
			Name:      "Time (s)",
			NameStyle: chart.StyleShow(),
			Style: chart.Style{
				Show: true,
			},
			ValueFormatter: func(v interface{}) string {
				if vf, isFloat := v.(float64); isFloat {
					// if vf == times[0] || vf == times[len(times)-1] {
					return fmt.Sprintf("%0.4f", vf)
					// }
				}
				return ""
			},
		},
		YAxis: chart.YAxis{
			Name:      "Cost",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		Series: charts,
	}

	collector := &chart.ImageWriter{}
	graph.Render(chart.PNG, collector)

	image, err := collector.Image()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Plotted Graph: %dx%d\n", image.Bounds().Size().X, image.Bounds().Size().Y)

	out, err := os.Create("output/output.jpg")
	defer out.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var opt jpeg.Options
	opt.Quality = 100

	// ok, write out the data into the new JPEG file
	err = jpeg.Encode(out, image, &opt) // put quality to 80%
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
