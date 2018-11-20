package main

import (
	"fmt"
	"math"
	"time"
)

// randomSearch finds a solution by finding random designs
// and comparing them against previous designs untill the
// the best is found.
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

// PSO (Particle Swarm Optimisation) finds the best design
// by creating a population of valid solutions, then flocking
// towards the best soloution in the population.
func PSO(noOfAntennae int, steeringAngle float64) (design []float64, peakSSL float64) {
	fmt.Printf("Starting PSO\n")

	// Init antennaArray
	var a AntennaArray
	a.noOfAntennae = noOfAntennae
	a.steeringAngle = steeringAngle

	var population []Particle

	// INITIALISE population, with random designs
	for i := 0; i < 30; i++ {
		currentParticle := Particle{}
		// Find a random design
		currentParticle.currentPostion = randomDesign(noOfAntennae)
		// Set the best peak values
		currentParticle.pBest = currentParticle.currentPostion
		currentParticle.pBestPeak, _ = a.evaluate(currentParticle.currentPostion)
		// Set the inital velocity to the differnce between the init position
		// and another random position divided by 2
		currentParticle.currentVelocity = make([]float64, noOfAntennae)
		tmp := randomDesign(noOfAntennae)
		for i, pos := range currentParticle.currentPostion {
			currentParticle.currentVelocity[i] = (tmp[i] - pos) / 2
		}

		population = append(population, currentParticle)
	}

	start := time.Now()

	var gBest []float64
	gBestPeak := math.MaxFloat64

	// Loop until time termination
	for i := 0; i >= 0; i++ {
		// Update global best
		for _, cParticle := range population {
			if cParticle.pBestPeak < gBestPeak {
				gBest = make([]float64, len(cParticle.currentPostion))
				copy(gBest, cParticle.currentPostion)
				gBestPeak = cParticle.pBestPeak
			}
		}

		// 1. UPDATE velocity and position
		// 2. EVALUATE new position
		// 3. UPDATE personal best
		for j := 0; j < len(population); j++ {
			population[j].update(gBest)

			// evaluate also updates the personal best
			population[j].evalulate(a)
		}

		// Termination condition
		now := time.Now()
		if now.Sub(start).Seconds() >= executeTime {
			fmt.Printf("Time Finshed!\n")
			break
		}
	}

	return gBest, gBestPeak
}
