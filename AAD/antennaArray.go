package main

import (
	"errors"
	"math"
	"sort"
)

const minSpacing = 0.25

// AntennaArray contains information such as the number of antennae
// and steeringAngel, as well as useful methods to valid designs
type AntennaArray struct {
	noOfAntennae  int
	steeringAngle float64
}

// bounds calculates the problem space min and max bounds
func (a *AntennaArray) bounds() (bnds [][]float64) {
	dimBnd := []float64{0, float64(a.noOfAntennae / 2.0)}
	for i := 0; i < a.noOfAntennae; i++ {
		bnds = append(bnds, dimBnd)
	}
	return
}

// isValid checks if the soloution is within the feasible region of the problem
// Constraints:
// 	A design is a vector of n_antennae anntena placements.
// 	A placement is a distance from the left hand side of the antenna array.
// 	A valid placement is one in which
// 		1) all antennae are separated by at least MIN_SPACING
//		2) the aperture size (the maximum element of the array) is exactly
//     		n_antennae/2.
func (a *AntennaArray) isValid(design []float64) bool {
	if len(design) != a.noOfAntennae {
		return false
	}
	sortedDesign := make([]float64, len(design))
	copy(sortedDesign, design)
	sort.Float64s(sortedDesign)

	// Check if size is noOfAntennae/2
	if math.Abs(sortedDesign[len(sortedDesign)-1]-float64(a.noOfAntennae)/2.0) > 1e-10 {
		return false
	}

	// Chekc if antennae are with the problem bounds
	for i := 0; i < len(sortedDesign)-1; i++ {
		if sortedDesign[i] < a.bounds()[i][0] ||
			sortedDesign[i] > a.bounds()[i][1] {
			return false
		}
	}

	// Check if antennae are correctly spaced
	for i := 0; i < len(sortedDesign)-1; i++ {
		if sortedDesign[i+1]-sortedDesign[i] < minSpacing {
			return false
		}
	}
	return true
}

// arrayFactor does some maths to work out the peak SSL.
func (a *AntennaArray) arrayFactor(design []float64, elevation float64) float64 {
	steering := 2.0 * math.Pi * a.steeringAngle / 360.0
	elevation = 2.0 * math.Pi * elevation / 360.0
	sum := 0.0
	for _, current := range design {
		sum += math.Cos(2 * math.Pi * current * (math.Cos(elevation) - math.Cos(steering)))
	}
	return 20.0 * math.Log(math.Abs(sum))
}

// evaluate give the peak SSL calculate from the design given.
// Designs which violate the constraints are penalised with a high cost.
func (a *AntennaArray) evaluate(design []float64) (peakSSL float64, err error) {
	if len(design) != a.noOfAntennae {
		return 0, errors.New("Design is the wrong size")
	}

	if !a.isValid(design) {
		// fmt.Printf("is not valid\n")
		return math.MaxFloat64, nil
	}

	type powerPeak struct {
		elevation float64
		power     float64
	}

	// Find all the peaks in power

	var peaks []powerPeak

	previous := powerPeak{0.0, math.SmallestNonzeroFloat64}
	current := powerPeak{0.0, a.arrayFactor(design, 0.0)}

	for elevation := 0.01; elevation <= 180.0; elevation += 0.01 {

		next := powerPeak{elevation, a.arrayFactor(design, elevation)}

		if current.power >= previous.power && current.power >= next.power {
			peaks = append(peaks, current)
		}
		previous = current
		current = next
	}
	peaks = append(peaks, powerPeak{180.0, a.arrayFactor(design, 180.0)})

	sort.Slice(peaks, func(i, j int) bool { return peaks[i].power > peaks[j].power })

	if len(peaks) < 2 {
		return math.SmallestNonzeroFloat64, nil
	}

	distanceFromSteering := math.Abs(peaks[0].elevation - a.steeringAngle)
	for i := 0; i < len(peaks); i++ {
		if math.Abs(peaks[i].elevation-a.steeringAngle) < distanceFromSteering {
			return peaks[0].power, nil
		}
	}
	return peaks[1].power, nil
}
