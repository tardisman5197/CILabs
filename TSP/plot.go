package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"

	chart "github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

func plot(costs []float64, times []float64) {
	fmt.Printf("Plotting Graph\n")
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
					if vf == times[0] || vf == times[len(times)-1] {
						return fmt.Sprintf("%0.4f", vf)
					}
				}
				return ""
			},
		},
		YAxis: chart.YAxis{
			Name:      "Cost",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: times,
				YValues: costs,
			},
		},
	}
	collector := &chart.ImageWriter{}
	graph.Render(chart.PNG, collector)

	image, err := collector.Image()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Plotted Graph: %dx%d\n", image.Bounds().Size().X, image.Bounds().Size().Y)

	out, err := os.Create("/output/output.jpg")
	defer out.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var opt jpeg.Options
	opt.Quality = 80

	// ok, write out the data into the new JPEG file
	err = jpeg.Encode(out, image, &opt) // put quality to 80%
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
