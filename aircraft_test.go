package airplaneslive

import (
	"testing"

	"github.com/fireflycons/geocoord"
)

func TestAircraftLocationReturnsCoordinateError(t *testing.T) {
	aircraft := Aircraft{Lat: 91, Lon: 0}
	if _, err := aircraft.Location(); err == nil {
		t.Fatal("expected invalid latitude error")
	}

	point := geocoord.MustNewCoordinate(0, 0)
	if _, err := aircraft.DistanceFrom(point); err == nil {
		t.Fatal("expected distance calculation to return invalid latitude error")
	}
	if _, err := aircraft.BearingFrom(point); err == nil {
		t.Fatal("expected bearing calculation to return invalid latitude error")
	}
}

func TestV2ResponseClosestToReturnsLocationError(t *testing.T) {
	response := V2Response{Ac: []Aircraft{{Lat: 91, Lon: 0}}}
	point := geocoord.MustNewCoordinate(0, 0)

	if _, err := response.ClosestTo(point); err == nil {
		t.Fatal("expected closest aircraft calculation to return invalid latitude error")
	}
	if _, err := response.ClosestToAirborne(point); err == nil {
		t.Fatal("expected closest airborne aircraft calculation to return invalid latitude error")
	}
}
