package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

const executeTime = 10

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// var a AntennaArray
	// a.noOfAntennae = 3
	// a.steeringAngle = 90
	// // design := []float64{0.5, 1.0, 1.5}
	// // peakSSL, err := a.evaluate(design)
	// // if err != nil {
	// // 	panic(err)
	// // }
	// // fmt.Printf("Peak SSL: %v\n", peakSSL)

	// design, peakSSL := randomSearch(3, 90)

	design, peakSSL := PSO(3, 90)
	fmt.Printf("Design: %v PeakSSL: %v\n", design, peakSSL)
}

func randomDesign(noOfAntennae int) (design []float64) {
	var a AntennaArray
	a.noOfAntennae = noOfAntennae

	for {
		var currentDesign []float64
		for i := 0; i < noOfAntennae-1; i++ {
			currentRand := 0.0 + rand.Float64()*((float64(noOfAntennae)/2)-0.0)
			currentDesign = append(currentDesign, currentRand)
		}
		currentDesign = append(currentDesign, float64(noOfAntennae)/2)

		sort.Float64s(currentDesign)

		if a.isValid(currentDesign) {
			design = currentDesign
			return
		}
	}
}
