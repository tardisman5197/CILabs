package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
)

const executeTime = 10

func main() {
	// demo()
	fmt.Printf("\n============\n")
	fmt.Printf("cities10.csv\n")
	cities := getCitiesFromFile("files/cities10.csv")

	// fmt.Printf("\n============\n")
	// fmt.Printf("cities16.csv\n")
	// cities = getCitiesFromFile("files/cities16.csv")

	fmt.Println("Finding Cheapest Route")
	route, cost, costs, times := localSearch(cities)
	fmt.Printf("Finished\nRoute: %v, Cost: %v\n", route, cost)
	plot(costs, times)
}

// demo finds the cheapest route from within the cities
// that are described by the demo given in the CI lab.
func demo() {
	cities := [][]float64{
		{0, 20, 42, 35},
		{20, 0, 30, 34},
		{42, 30, 0, 12},
		{35, 34, 12, 0},
	}

	fmt.Printf("Cities: %v\n\n", cities)

	// // Generate a random amount of routes
	// fmt.Printf("Random Routes:\n")
	// for i := 0; i < 10; i++ {
	// 	route := generateRandomRoute(len(cities))

	// 	cost := getCostOfRoute(cities, route)

	// 	fmt.Printf("\t%v - Route: %v, Cost: %v\n", i, route, cost)
	// }

	// randomSearch(cities)

	route, cost, _, _ := localSearch(cities)
	fmt.Printf("Cheapest Route: %v - %v\n", route, cost)

}

// getCostOfRoute takes in a route and outputs the distance.
// The cost is calculated by the weights between the cities and
// includes the cost to get back to the starting city.
func getCostOfRoute(cities [][]float64, route []int) (cost float64) {
	previousCity := -1
	for _, city := range route {
		if previousCity > -1 {
			cost += cities[city][previousCity]
		}
		previousCity = city
	}
	cost += cities[0][len(cities)-1]
	return
}

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

// twoOpt finds all of the permutations when swapping two cities in a route.
// This is done by itterating through each possibilities (e.g. 0 <-> 0, 0 <-> 1 ...)
// and only adding routes that have not been found previous to the list of routes.
func twoOpt(route []int) (twoOptRoutes [][]int) {
	usedRoutes := make(map[string]bool)
	for i := 0; i < len(route); i++ {
		for j := 0; j < len(route); j++ {
			currentRoute := make([]int, len(route))
			copy(currentRoute, route)

			// Swap i and j
			tmp := currentRoute[i]
			currentRoute[i] = currentRoute[j]
			currentRoute[j] = tmp

			// Check if route has already been used
			currentRouteStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(currentRoute)), ""), "[]")
			if _, used := usedRoutes[currentRouteStr]; !used {
				usedRoutes[currentRouteStr] = true
				twoOptRoutes = append(twoOptRoutes, currentRoute)
			}
		}
	}
	return
}

// bestNeighbourStep findes the cheapest route out of a set of routes.
func bestNeighbourStep(cities [][]float64, routes [][]int) (route []int) {
	cheapest := math.MaxFloat64
	for _, currentRoute := range routes {
		cost := getCostOfRoute(cities, currentRoute)
		// fmt.Printf("%v - %v\n", currentRoute, cost)
		if cost < cheapest {
			cheapest = cost
			route = currentRoute
		}
	}
	return
}
