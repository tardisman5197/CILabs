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
func randomSearch(cities [][]float64) (route []int, cost float64, line [][]float64) {
	start := time.Now()
	line = [][]float64{{}, {}}

	cost = math.MaxFloat64

	var noOfPermutations big.Int
	noOfPermutations.MulRange(1, int64(len(cities)-1))

	// fmt.Printf("Number of Permutations: %v\n", noOfPermutations.Int64())

	usedRoutes := make(map[string]bool)

	// Loop through all permutations
	for i := int64(0); i < noOfPermutations.Int64(); i++ {

		// if i%10000 == 0 {
		// 	percentageChecked := math.Ceil(float64(i+1) / float64(noOfPermutations.Int64()) * 100)
		// 	// fmt.Printf("Pervcentage = %v\n", percentageChecked)
		// 	// fmt.Printf("\rProgress: %v/%v - %v%%", i+1, noOfPermutations.Int64(), percentageChecked)
		// }

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

					now := time.Now()
					line[1] = append(line[1], cost)
					line[0] = append(line[0], now.Sub(start).Seconds())

					fmt.Printf("New cheapest found - %v : %v\n", route, cost)
				}
				break
			}

			// Check if ran out of time
			now := time.Now()
			if now.Sub(start).Seconds() >= executeTime {
				break
			}
		}

		// Check if ran out of time
		now := time.Now()
		if now.Sub(start).Seconds() >= executeTime {
			fmt.Printf("Execute Time Acheived\n")
			now := time.Now()
			line[1] = append(line[1], cost)
			line[0] = append(line[0], now.Sub(start).Seconds())
			break
		}
	}
	return
}

func localSearch(cities [][]float64) (globalBestRoute []int, globalBestCost float64, line [][]float64) {
	start := time.Now()
	line = [][]float64{{}, {}}

	globalBestCost = math.MaxFloat64

	// Start a local search.
	for {
		// Generate starting route
		localBestRoute := generateRandomRoute(len(cities))
		localBestCost := getCostOfRoute(cities, localBestRoute)

		// Find local optimal value
		for {
			neighbourhood := twoOpt(localBestRoute)

			currentRoute, currentCost := bestNeighbourStep(cities, neighbourhood)
			// fmt.Printf("%v - %v\n", currentRoute, currentCost)

			if currentCost < localBestCost {
				// New cheapest found
				localBestRoute = currentRoute
				localBestCost = currentCost

			} else if currentCost == localBestCost {
				// No more improvements being made
				localBestRoute = currentRoute
				localBestCost = currentCost
				break
			}

			// Check if ran out of time
			now := time.Now()
			if now.Sub(start).Seconds() >= executeTime {
				break
			}
		}

		// Check if this local search has beaten the best overall
		if localBestCost < globalBestCost {
			globalBestRoute = localBestRoute
			globalBestCost = localBestCost

			fmt.Printf("Now Cheapest Route Found: %v - %v\n", globalBestRoute, globalBestCost)

			// Add new best to the graph points
			now := time.Now()
			line[1] = append(line[1], globalBestCost)
			line[0] = append(line[0], now.Sub(start).Seconds())
		}

		// Check if ran out of time
		now := time.Now()
		if now.Sub(start).Seconds() >= executeTime {
			fmt.Printf("Execute Time Acheived\n")

			now := time.Now()
			line[1] = append(line[1], globalBestCost)
			line[0] = append(line[0], now.Sub(start).Seconds())
			break
		}
	}
	return
}
