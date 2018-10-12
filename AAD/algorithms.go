package main

import (
	"math"
	"time"
)

func randomSearch(noOfAntennae int, steeringAngle float64) (design []float64, peakSSL float64) {
	start := time.Now()

	peakSSL = math.MaxFloat64

	var a AntennaArray
	a.noOfAntennae = noOfAntennae
	a.steeringAngle = steeringAngle

	for {
		currentDesign := randomDesign(noOfAntennae)

		currentPeakSSL, err := a.evaluate(currentDesign)
		if err != nil {
			panic(err)
		}

		if currentPeakSSL < peakSSL {
			peakSSL = currentPeakSSL
			design = currentDesign
		}

		now := time.Now()
		if now.Sub(start).Seconds() >= executeTime {
			break
		}
	}
	return
}
