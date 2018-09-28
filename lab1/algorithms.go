package main

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"strings"
)

// generateRandomRouteHat takes in the number of cities and outputs a random valid route.
// Routes are made by palcing all the cities inside a hat. Then randomly removes
// a city from the hat and addes it to the route.
func generateRandomRoute(noOfCities int) (route []int) {
	var hat []int

	// Always start with the last city
	route = []int{noOfCities - 1}

	// Fill the hat with cities
	for i := 0; i < noOfCities-1; i++ {
		hat = append(hat, i)
	}

	for i := 0; i < noOfCities-1; i++ {
		city := rand.Intn(len(hat))

		route = append(route, hat[city])

		// remove city from the hat
		hat[city] = hat[len(hat)-1]
		hat = hat[:len(hat)-1]
	}

	return
}

// randomSearch finds the cheapest route around all the cities.
// It randomly generates a route, checks if it has been seen before,
// then checks if the route found is cheaper than any previous routes.
// Once all the permutations have been checked it returns the cheapest
// route.
func randomSearch(cities [][]float64) (route []int, cost float64) {
	cost = math.MaxFloat64

	var noOfPermutations big.Int
	noOfPermutations.MulRange(1, int64(len(cities)-1))

	fmt.Printf("Number of Permutations: %v\n", noOfPermutations.Int64())

	usedRoutes := make(map[string]bool)

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
					fmt.Printf("\nNew cheapest found - %v : %v\n", route, cost)
				}
				break
			}
		}
	}
	return
}
