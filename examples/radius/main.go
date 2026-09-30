package main

import (
	"context"
	"fmt"
	"log"

	"github.com/fireflycons/airplaneslive"
	"github.com/fireflycons/geocoord"
)

func main() {
	ctx := context.Background()
	api, err := airplaneslive.NewApi()
	if err != nil {
		log.Fatalf("Failed to create API instance: %v", err)
	}

	// London Centre - generally regarded as Charing Cross
	london := geocoord.MustNewCoordinate(51.5074, -0.1278) // London coordinates
	radius := 10.0

	response, err := api.AircraftWithinRadius(ctx, london, radius)
	if err != nil {
		log.Fatalf("Failed to fetch aircraft data: %v", err)
	}

	fmt.Printf("%d Aircraft within radius\n", len(response.Ac))

	if closestPlane := response.ClosestTo(london); closestPlane != nil {

		fmt.Printf(
			"Closest plane to central London: %s (type: %s, heading: %.1f, alt: %d, climb: %s, bank: %s, dist: %.1f, bearing: %.1f)",
			closestPlane.Flight,
			closestPlane.TypeCode,
			closestPlane.TrueHeading,
			closestPlane.AltBaro.Altitude,
			closestPlane.Pitch(),
			closestPlane.Bank(),
			closestPlane.DistanceFrom(london),
			closestPlane.BearingFrom(london),
		)
	}
}
