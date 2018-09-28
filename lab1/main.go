package main

import (
	"fmt"
)

func main() {
	demo()
	fmt.Printf("\n============\n")

	cities := getCitiesFromFile("ulysses16.csv")

	fmt.Println("Finding Cheapest Route")
	randomSearch(cities)
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

	// Generate a random amount of routes
	fmt.Printf("Random Routes:\n")
	for i := 0; i < 10; i++ {
		route := generateRandomRoute(len(cities))

		cost := getCostOfRoute(cities, route)

		fmt.Printf("\t%v - Route: %v, Cost: %v\n", i, route, cost)
	}

	randomSearch(cities)
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
