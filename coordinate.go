package airplaneslive

import (
	"fmt"
	"math"
)

// Coordinate represents a geographic position in decimal degrees.
type Coordinate struct {
	lat float64
	lon float64
}

// NewCoordinate creates a coordinate and returns an error if either value is out of range.
func NewCoordinate(lat, lon float64) (Coordinate, error) {
	if err := checkLat(lat); err != nil {
		return Coordinate{}, err
	}
	if err := checkLon(lon); err != nil {
		return Coordinate{}, err
	}
	return Coordinate{lat: lat, lon: lon}, nil
}

// MustNewCoordinate creates a coordinate and panics if either value is out of range.
func MustNewCoordinate(lat, lon float64) Coordinate {
	coord, err := NewCoordinate(lat, lon)
	if err != nil {
		panic(fmt.Sprintf("invalid coordinate: %v", err))
	}
	return coord
}

// DistanceTo calculates the distance in nautical miles from this coordinate to other.
func (co Coordinate) DistanceTo(other Coordinate) float64 {

	// Earth's mean radius in nautical miles (6371.0088 km / 1.852 km/NM)
	const earthRadiusNM = 3440.0695

	// Convert latitudes and longitudes from degrees to radians
	lat1 := co.lat * math.Pi / 180
	lon1 := co.lon * math.Pi / 180
	lat2 := other.lat * math.Pi / 180
	lon2 := other.lon * math.Pi / 180

	// Haversine formula
	dLat := lat2 - lat1
	dLon := lon2 - lon1

	sinDLat := math.Sin(dLat / 2)
	sinDLon := math.Sin(dLon / 2)

	h := (sinDLat * sinDLat) + math.Cos(lat1)*math.Cos(lat2)*(sinDLon*sinDLon)
	c := 2 * math.Asin(math.Min(1, math.Sqrt(h))) // math.Min prevents rounding issues exceeding 1.0

	return earthRadiusNM * c
}

// HeadingTo calculates the initial bearing when starting at this point to get to the destination coordinate.
func (co Coordinate) HeadingTo(other Coordinate) float64 {

	// Convert latitudes and longitudes from degrees to radians
	lat1 := co.lat * math.Pi / 180
	lon1 := co.lon * math.Pi / 180
	lat2 := other.lat * math.Pi / 180
	lon2 := other.lon * math.Pi / 180

	dLon := lon2 - lon1

	// Formula for initial bearing
	y := math.Sin(dLon) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(dLon)

	// math.Atan2 returns values from -Pi to +Pi radians
	headingRad := math.Atan2(y, x)

	// Convert radians back to degrees
	headingDeg := headingRad * 180 / math.Pi

	// Normalize heading to standard compass degrees (0° to 360°)
	// Go's math.Mod handles negative values gracefully this way
	return math.Mod(headingDeg+360.0, 360.0)
}

func checkLat(lat float64) error {
	if lat < -90 || lat > 90 {
		return fmt.Errorf("invalid latitude %f: must be in [-90, 90]", lat)
	}
	return nil
}

func checkLon(lon float64) error {
	if lon < -180 || lon > 180 {
		return fmt.Errorf("invalid longitude %f: must be in [-180, 180]", lon)
	}
	return nil
}
