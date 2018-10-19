package main

import (
	"fmt"
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

// PSO ..
func PSO(noOfAntennae int, steeringAngle float64) (design []float64, peakSSL float64) {
	var a AntennaArray
	a.noOfAntennae = noOfAntennae
	a.steeringAngle = steeringAngle

	var population []Particle
	// INITIALISE population
	fmt.Printf("Init Population\n")
	for i := 0; i < 30; i++ {
		currentParticle := Particle{}
		currentParticle.currentPostion = randomDesign(noOfAntennae)
		currentParticle.pBest = currentParticle.currentPostion
		currentParticle.pBestPeak, _ = a.evaluate(currentParticle.currentPostion)
		currentParticle.currentVelocity = make([]float64, noOfAntennae)
		tmp := randomDesign(noOfAntennae)
		for i, pos := range currentParticle.currentPostion {
			currentParticle.currentVelocity[i] = (tmp[i] - pos) / 2
		}
		population = append(population, currentParticle)
	}
	// fmt.Printf("Population:\n%v\n", population)

	start := time.Now()

	var gBest []float64
	gBestPeak := math.MaxFloat64

	fmt.Printf("Start iterations\n")
	for i := 0; i >= 0; i++ {
		// Update global best
		for _, cParticle := range population {
			if cParticle.pBestPeak < gBestPeak {
				gBest = make([]float64, len(cParticle.currentPostion))
				copy(gBest, cParticle.currentPostion)
				gBestPeak = cParticle.pBestPeak
				fmt.Printf("%v: New gBest: %v : %v\n", i, gBest, gBestPeak)
			}
		}

		// 1. UPDATE velocity and position
		// 2. EVALUATE new position
		// 3. UPDATE personal best
		for j := 0; j < len(population); j++ {
			population[j].update(gBest)
			// evaluate also updates the personal best
			// if j == 0 {
			// 	fmt.Printf("%v: %v : V %v \n%v: %v : %v\n", i, population[j].currentPostion, population[j].currentVelocity, i, population[j].pBest, population[j].pBestPeak)
			// }
			population[j].evalulate(a)
			// if j == 0 {
			// 	fmt.Printf("Better: %v\n", better)
			// 	fmt.Printf("%v: %v : V %v \n%v: %v : %v\n\n", i, population[j].currentPostion, population[j].currentVelocity, i, population[j].pBest, population[j].pBestPeak)
			// }

		}

		// Termination condition
		now := time.Now()
		if now.Sub(start).Seconds() >= executeTime {
			fmt.Printf("Time Finshed\n")
			break
		}
	}

	return gBest, gBestPeak
}
