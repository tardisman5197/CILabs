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
	// task2()
}

func task1() {
	fmt.Println("Lab 1")

	cities := [][]float64{
		{0, 22, 42, 35},
		{20, 0, 30, 34},
		{42, 30, 0, 12},
		{35, 34, 12, 0},
	}

	fmt.Printf("Cities: %v\n", cities)

	// Generate a random amount of routes
	// for i := 0; i < 10; i++ {
	// 	route := generateRandomRoute(cities)

	// 	cost := getCostOfRoute(cities, route)

	// 	fmt.Printf("%v - Route: %v, Cost: %v\n", i, route, cost)
	// }

	routes := getPermutations(len(cities))

	route, cost := getCheapestRoute(cities, routes)
	fmt.Printf("Cheapest Route: %v, Cost: %v\n", route, cost)
}

func task2() {
	coords := readCSVFile("ulysses16.csv")
	cities := getCitiesFromCoords(coords)
	for i, city := range cities {
		fmt.Printf("ID: %v, Connections:\n", i)
		for j, distance := range city {
			fmt.Printf("\t%v -> %v = %v\n", i, j, distance)
		}
	}
	// route := generateRandomRouteHat(len(cities))
	// cost := getCostOfRoute(cities, route)
	// fmt.Printf("Route: %v, Costs: %v\n", route, cost)
	routes := getPermutations(len(cities))
	fmt.Printf("Routes: %v\n", routes)
}

func validCity(route []int, city int) bool {
	for _, currentCity := range route {
		if currentCity == city {
			return false
		}
	}
	return true
}

func validRoute(cities [][]float64, route []int) bool {
	if len(route) != len(cities) {
		return false
	}

	check := make(map[int]bool)
	for _, city := range route {
		if _, used := check[city]; used || city >= len(cities) {
			return false
		} else {
			check[city] = true
		}
	}
	return true
}

func getCostOfRoute(cities [][]float64, route []int) (cost float64) {
	previousCity := -1
	for _, city := range route {
		if previousCity > -1 {
			cost += cities[city][previousCity]
		}
		previousCity = city
	}
	return
}

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
func generateRandomRoute(cities [][]float64) (route []int) {
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

func generateRandomRouteHat(noOfCities int) (route []int) {
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

func getPermutations(noOfCities int) (routes [][]int) {

	var noOfPermutations big.Int
	noOfPermutations.MulRange(1, int64(noOfCities-1))

	fmt.Printf("No of permutations: %v\n", noOfPermutations.Int64())
	usedRoutes := make(map[string]bool)

	for i := int64(0); i < noOfPermutations.Int64(); i++ {
		var currentRoute []int
		for {
			currentRoute = generateRandomRouteHat(noOfCities)
			currentRouteStr := convertToString(currentRoute)
			if _, used := usedRoutes[currentRouteStr]; !used {
				usedRoutes[currentRouteStr] = true
				break
			}
		}
		routes = append(routes, currentRoute)
	}
	return
}

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

func convertToString(a []int) (s string) {
	for _, val := range a {
		s += strconv.Itoa(val) + ","
	}
	return
}
