package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"math/big"
	"math/rand"
	"os"
	"strconv"
)

func main() {
	task1()
	fmt.Printf("\n============\n")
	// task2()
}

// Steps 1-4
func task1() {
	fmt.Println("Tasks 1-4")

	cities := [][]float64{
		{0, 22, 42, 35},
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

	// routes := getPermutations(len(cities))
	// fmt.Printf("\nNo of permutations: %v\n", len(routes))

	// route, cost := getCheapestRoute(cities, routes)
	// fmt.Printf("Cheapest Route: %v, Cost: %v\n", route, cost)

	randomSearch(cities)
}

// Steps 5-6
func task2() {
	fmt.Println("Tasks 5-6")

	coords := readCSVFile("ulysses16.csv")
	cities := getCitiesFromCoords(coords)

	for i, city := range cities {
		fmt.Printf("ID: %v, Connections:\n", i)
		for j, distance := range city {
			fmt.Printf("\t%v -> %v = %v\n", i, j, distance)
		}
	}

	routes := getPermutations(len(cities))
	fmt.Printf("Routes: %v\n", routes)
}

// validCity checks if the city provided can be added to the route
func validCity(route []int, city int) bool {
	for _, currentCity := range route {
		if currentCity == city {
			return false
		}
	}
	return true
}

// validRoute checks if the route if the route is alowd
func validRoute(cities [][]float64, route []int) bool {
	if len(route) != len(cities) {
		return false
	}

	check := make(map[int]bool)
	for _, city := range route {
		if _, used := check[city]; used || city >= len(cities) {
			return false
		}
		check[city] = true
	}
	return true
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

// getCheapestRoute takes in a bunch of routes and outputs the cheapest route and its cost
func getCheapestRoute(cities [][]float64, routes [][]int) (route []int, cost float64) {
	cost = math.MaxFloat64
	for _, currentRoute := range routes {
		currentCost := getCostOfRoute(cities, currentRoute)
		if currentCost < cost {
			cost = currentCost
			route = currentRoute
		}
	}
	return
}

// generateRandomRoute takes in a list of cities and outputs a random valid route
func generateRandomRouteOld(cities [][]float64) (route []int) {
	for i := 0; i < len(cities); i++ {
		for {
			city := rand.Intn(len(cities))
			if validCity(route, city) {
				route = append(route, city)
				break
			}
		}
	}
	return
}

// generateRandomRouteHat takes in the number of cities and outputs a random valid route.
// Routes are made my palcing all the cities inside a hat. Then randomly removes
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

// getPermutations gets all possible routes between all the cities
func getPermutations(noOfCities int) (routes [][]int) {

	var noOfPermutations big.Int
	noOfPermutations.MulRange(1, int64(noOfCities-1))

	usedRoutes := make(map[string]bool)

	for i := int64(0); i < noOfPermutations.Int64(); i++ {
		var currentRoute []int
		for {
			currentRoute = generateRandomRoute(noOfCities)
			currentRouteStr := arrayToString(currentRoute)
			if _, used := usedRoutes[currentRouteStr]; !used {
				usedRoutes[currentRouteStr] = true
				break
			}
		}
		routes = append(routes, currentRoute)
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

		var currentRoute []int

		// Find a random route that has not been found
		for {
			currentRoute = generateRandomRoute(len(cities))
			currentRouteStr := arrayToString(currentRoute)

			// Check if the route has already been found
			if _, used := usedRoutes[currentRouteStr]; !used {
				usedRoutes[currentRouteStr] = true

				// Check if the route is cheaper than the current cheapest
				currentCost := getCostOfRoute(cities, currentRoute)
				if currentCost < cost {
					route = currentRoute
					cost = currentCost
					fmt.Printf("New cheapest found - %v : %v", route, cost)
				}
				break
			}
		}
	}
	return
}

// getCitiesFromFile reads a csv file and converts the coordinates
// read into distances between cities.
func getCitiesFromFile(path string) (cities [][]float64) {
	coords := readCSVFile(path)
	cities = getCitiesFromCoords(coords)
	return
}

// readCSVFile reads a file and outputs a list of the coords for each city
func readCSVFile(path string) (coords [][]float64) {
	csvFile, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	reader.Read()
	reader.Read()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		x, _ := strconv.ParseFloat(record[1], 64)

		y, _ := strconv.ParseFloat(record[2], 64)
		coords = append(coords, []float64{x, y})
	}
	fmt.Println(coords)

	return
}

// getCitiesFromCoords converts a list of coords to a list of cities and distances
func getCitiesFromCoords(coords [][]float64) (cities [][]float64) {
	for _, source := range coords {
		var distances []float64
		for _, destination := range coords {
			distance := math.Sqrt(
				(destination[0]-source[0])*(destination[0]-source[0]) +
					(destination[1]-source[1])*(destination[1]-source[1]))
			distances = append(distances, distance)
		}
		cities = append(cities, distances)
	}
	return
}

// arrayToString converts and array of ints to a string with commas
func arrayToString(a []int) (s string) {
	for _, val := range a {
		s += strconv.Itoa(val) + ","
	}
	return
}
