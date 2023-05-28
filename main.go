package main

import (
	"math"
	"os"
	"sort"

	"github.com/gocarina/gocsv"
)

const radConv = math.Pi / 180.0

type Station struct {
	Name      string  `csv:"Station"`
	Latitude  float64 `csv:"Latitude"`
	Longitude float64 `csv:"Longitude"`
}

type OutputStation struct {
	Name     string  `csv:"station_name"`
	Distance float64 `csv:"distance"`
}

type Position struct {
	Latitude  float64
	Longitude float64
}

func main() {
	// Read csv
	bytes, err := os.ReadFile("london_stations.csv")
	if err != nil {
		panic(err)
	}

	var stations []Station
	err = gocsv.UnmarshalBytes(bytes, &stations)
	if err != nil {
		panic(err)
	}

	target := Position{Latitude: 51.479495, Longitude: -0.000500}

	outputStations := make([]OutputStation, len(stations))
	for i, station := range stations {
		outputStations[i] = OutputStation{
			Name:     station.Name,
			Distance: distance(target, Position{Latitude: station.Latitude, Longitude: station.Longitude}),
		}
	}

	// sort stations by distance
	sort.SliceStable(outputStations, func(i, j int) bool {
		return outputStations[i].Distance < outputStations[j].Distance
	})

	file, err := os.Create("closest_stations.csv")
	if err != nil {
		panic(err)
	}
	gocsv.MarshalFile(&outputStations, file)
}

func distance(p, q Position) float64 {
	dlong := (q.Longitude - p.Latitude) * radConv
	dlat := (q.Latitude - p.Latitude) * radConv
	a := math.Pow(math.Sin(dlat/2.0), 2) + math.Cos(p.Latitude*radConv)*math.Cos(q.Latitude*radConv)*math.Pow(math.Sin(dlong/2.0), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return 6357 * c
}
