package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
)

// getCitiesFromFile reads a csv file and converts the coordinates
// read into distances between cities.
func getCitiesFromFile(path string) (cities [][]float64) {
	coords := readCSVFile(path)
	cities = getCitiesFromCoords(coords)
	return
}

// readCSVFile reads a file and outputs a list of the coords for each city.
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
