package main

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"sort"
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

// localSearch finds the best route it can in a given time. This is
// done by generating a random route, then getting all the variations
// of that route. From them permitations the best is found to start the
// next iteration. This repeats until the best in that nbourhood is found.
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

// evolutionaryAlgorithm finds a solution by applying methods used by evolution.
func evolutionaryAlgorithm(cities [][]float64) (bestRoute []int, bestCost float64, line [][]float64) {
	var population []Route

	// Init population with random routes and Eval the routes.
	fmt.Printf("Init pop\n")
	for i := 0; i < populationSize; i++ {
		var currentRoute Route
		currentRoute.route = generateRandomRoute(len(cities))
		currentRoute.cost = getCostOfRoute(cities, currentRoute.route)
		population = append(population, currentRoute)
	}

	fmt.Printf("Population: %v\n", population)

	// Repeat until termination

	for i := 0; i < 10; i++ {
		// Parent Selection
		// Tournament slection:
		//	1. Pick k memebers at random
		//	2. Choose the best out of the selection
		//	3. Repeate until pop size reached
		k := 2
		var parents []Route

		fmt.Printf("%v: Select Parents\n", i)
		for j := 0; j < populationSize; j++ {
			var parent Route
			parent.cost = math.MaxFloat64
			for m := 0; m < k; m++ {
				// Get next random parent
				currentParent := population[rand.Intn(populationSize)]
				// Check if better then the current best
				if currentParent.cost < parent.cost {
					parent = currentParent
				}
			}
			parents = append(parents, parent)
			fmt.Printf("\t%v: %v\n", j, parent)
		}

		var offspring []Route

		fmt.Printf("%v: Create offspring\n", i)
		for j := 0; j < populationSize; j++ {
			fmt.Printf("\t%v:\n", j)
			var currentOffspring Route
			// Recombine parents
			currentOffspring.route = orderOneCrossover(parents[j].route, parents[rand.Intn(populationSize)].route)
			fmt.Printf("\t\tComb: %v\n", currentOffspring)
			// Mutate
			if rand.Float64() <= mutateProbability {
				// Swap two cities at random
				x := rand.Intn(len(currentOffspring.route))
				y := rand.Intn(len(currentOffspring.route))
				tmp := currentOffspring.route[x]
				currentOffspring.route[x] = currentOffspring.route[y]
				currentOffspring.route[y] = currentOffspring.route[tmp]
			}
			fmt.Printf("\t\tMuta: %v\n", currentOffspring)
			// Evaluate
			currentOffspring.cost = getCostOfRoute(cities, currentOffspring.route)
			fmt.Printf("\t\tEval: %v\n", currentOffspring)
			offspring = append(offspring, currentOffspring)
		}

		// Select next gen
		population = offspring
	}

	bestCost = math.MaxFloat64
	for i := 0; i < len(population); i++ {
		if population[i].cost < bestCost {
			bestRoute = population[i].route
			bestCost = population[i].cost
		}
	}
	return
}

// artificialImmuneSystem finds a solution for TSP, by using methods similar to
// an immune system. The steps that this algorithm takes are as follows:
// 	1. Initiation, create random soloutions
//	2. Cloning, make beta amount of copies
//	3. Mutation, inverse proportional hyper-mutation
//	4. Selection, choose the best mu for the next population
//	5. Metadynamics, repace the worst d with random solutions
//	6. Repeat until termination condition
func artificialImmuneSystem(cities [][]float64) (bestRoute []int, bestCost float64, line [][]float64) {
	var population []Route

	// Init population with random routes and Eval the routes.
	population = generateRandomPopulation(cities, populationSize)

	// Repeat until terminating condition
	for i := 0; i < 5000; i++ {
		// Cloning
		var clonePool []Route
		for j := 0; j < len(population); j++ {
			for k := 0; k < cloneSizeFactor; k++ {
				currentClone := Route{}
				currentClone.cost = population[j].cost
				currentClone.route = make([]int, len(population[j].route))
				copy(currentClone.route, population[j].route)
				clonePool = append(clonePool, currentClone)
			}
		}

		// Mutation
		for j := 0; j < len(clonePool); j++ {
			// Random hotspot
			start := rand.Intn(len(clonePool[j].route))

			size := len(clonePool[j].route)

			// length = routeLength * exp(-p*f/fBest)
			inv := math.Exp(-0.5 * (clonePool[j].cost / bestFitness))
			lengthFloat := inv * float64(size)
			length := int(lengthFloat)

			// Reverse the section
			var tmp []int
			for k := 0; k < length; k++ {
				tmp = append(tmp, clonePool[j].route[(k+start)%size])
			}

			for k := 0; k < length; k++ {
				clonePool[j].route[(k+start)%size] = tmp[len(tmp)-(k+1)]
			}

			clonePool[j].cost = getCostOfRoute(cities, clonePool[j].route)
		}

		// Selection
		population = append(population, clonePool...)

		sort.SliceStable(population, func(i, j int) bool { return population[i].cost < population[j].cost })

		population = population[:populationSize]

		// Metadynamics
		for j := 1; j <= replacementSize; j++ {
			var currentRoute Route
			currentRoute.route = generateRandomRoute(len(cities))
			currentRoute.cost = getCostOfRoute(cities, currentRoute.route)
			population[len(population)-j] = currentRoute
		}

		bestRoute = population[0].route
		bestCost = population[0].cost
	}
	return
}
