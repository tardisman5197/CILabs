package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

const executeTime = 5
const populationSize = 10
const mutateProbability = 0.7

func main() {
	// demo()
	rand.Seed(time.Now().UTC().UnixNano())

	// fmt.Printf("\n============\n")
	// fmt.Printf("cities10.csv\n")
	// cities := getCitiesFromFile("files/cities10.csv")

	// fmt.Printf("\n============\n")
	// fmt.Printf("cities16.csv\n")
	// cities := getCitiesFromFile("files/cities16.csv")

	fmt.Println("Finding Cheapest Route")
	// route, cost, randomLine := randomSearch(cities)
	// fmt.Printf("Random Finished\nRoute: %v, Cost: %v\n", route, cost)
	// route, cost, localLine := localSearch(cities)
	// fmt.Printf("Local Finished\nRoute: %v, Cost: %v\n", route, cost)
	// plot(randomLine, localLine)

	// route, cost, _ := evolutionaryAlgorithm(cities)
	// fmt.Printf("Evolution Finished\nRoute: %v, Cost: %v\n", route, cost)

	p1 := []int{0, 1, 2, 3, 4, 5}
	p2 := []int{5, 4, 3, 2, 1, 0}
	fmt.Printf("Child: %v\n", orderOneCrossover(p1, p2))
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

	route, cost, _ := localSearch(cities)
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
func bestNeighbourStep(cities [][]float64, routes [][]int) (route []int, cost float64) {
	cost = math.MaxFloat64
	for _, currentRoute := range routes {
		currentCost := getCostOfRoute(cities, currentRoute)
		// fmt.Printf("%v - %v\n", currentRoute, cost)
		if currentCost < cost {
			cost = currentCost
			route = currentRoute
		}
	}
	return
}

// orderOneCrossover takes two parents and combines them to make a child.
// This is achived by random selecting a section of one parent to make up part
// of the offspring then filling the remaining gaps with the cities not used in
// the order of the other parent.
func orderOneCrossover(p1 []int, p2 []int) (child []int) {
	child = make([]int, len(p1))

	// Randomly select the section of one parent and give to child
	start := rand.Intn(len(p1))
	end := rand.Intn(len(p1))

	if end < start {
		tmp := start
		start = end
		end = tmp
	}
	fmt.Printf("start: %v\n", start)
	fmt.Printf("end: %v\n", end)

	used := make(map[int]bool)
	for i := start; i < end; i++ {
		child[i] = p1[i]
		used[p1[i]] = true
	}

	// Fill gaps in child with other parent
	var tmp []int
	for i := 0; i < len(p2); i++ {
		if _, ok := used[p2[i]]; !ok {
			tmp = append(tmp, p2[i])
		}
	}

	for i := 0; i < len(child); i++ {
		if i >= start && i < end {
			// Get next city
			child[i] = tmp[0]
			// pop city from p2
			if len(tmp) > 1 {
				tmp = tmp[1:]
			}
		}
	}

	return
}
