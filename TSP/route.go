package main

// Route stores a solution for TSP. This is used by the
// evolutionary algorithm.
type Route struct {
	route []int
	cost  float64
}
