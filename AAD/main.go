package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

const executeTime = 10

const numberofAntennae = 5
const steeringAngle = 90

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	fmt.Printf("Execute time: %v\nNumber Of Antennae: %v\nSteering Angle: %v\n\n", executeTime, numberofAntennae, steeringAngle)
	design, peakSSL := PSO(numberofAntennae, steeringAngle)
	fmt.Printf("Design: %v PeakSSL: %v\n", design, peakSSL)
}

// randomDesign gets a random valid design.
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
