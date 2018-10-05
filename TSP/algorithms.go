package main

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"
)

// randomSearch finds the cheapest route around all the cities.
// It randomly generates a route, checks if it has been seen before,
// then checks if the route found is cheaper than any previous routes.
// Once all the permutations have been checked it returns the cheapest
// route.
func randomSearch(cities [][]float64) (route []int, cost float64, costs []float64, times []float64) {
	cost = math.MaxFloat64

	var noOfPermutations big.Int
	noOfPermutations.MulRange(1, int64(len(cities)-1))

	fmt.Printf("Number of Permutations: %v\n", noOfPermutations.Int64())

	usedRoutes := make(map[string]bool)

	start := time.Now()
	// Loop through all permutations
	for i := int64(0); i < noOfPermutations.Int64(); i++ {

		if i%10000 == 0 {
			percentageChecked := math.Ceil(float64(i+1) / float64(noOfPermutations.Int64()) * 100)
			// fmt.Printf("Pervcentage = %v\n", percentageChecked)
			fmt.Printf("\rProgress: %v/%v - %v%%", i+1, noOfPermutations.Int64(), percentageChecked)
		}

		var currentRoute []int

		// Find a random route that has not been found
		for {
			currentRoute = generateRandomRoute(len(cities))
			currentRouteStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(currentRoute)), ""), "[]")

			// Check if the route has already been found
			if _, used := usedRoutes[currentRouteStr]; !used {
				usedRoutes[currentRouteStr] = true

				// Check if the route is cheaper than the current cheapest
				currentCost := getCostOfRoute(cities, currentRoute)
				if currentCost < cost {
					route = currentRoute
					cost = currentCost
					costs = append(costs, cost)
					now := time.Now()
					times = append(times, now.Sub(start).Seconds())
					fmt.Printf("\nNew cheapest found - %v : %v\n", route, cost)
				}
				break
			}
		}
	}
	return
}

func localSearch(cities [][]float64) (route []int, cost float64, costs []float64, times []float64) {
	return
}
