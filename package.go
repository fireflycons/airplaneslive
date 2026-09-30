// Package airplaneslive provides a Go client for aircraft position and
// reference-data queries against the Airplanes.live API.
//
// Create an [Api] with [NewApi], then pass a context to each request. Aircraft
// queries include positions within a radius and positions by registration.
// Reference-data methods search or list airlines, airports, cities, countries,
// and timezones.
//
// Coordinates use [geocoord.Coordinate] values from github.com/fireflycons/geocoord.
// Its NewCoordinate and MustNewCoordinate constructors validate latitude and
// longitude ranges. Aircraft distances and radius values are in nautical miles,
// and headings are in degrees.
//
// Reference-data results are paginated. Pass [DefaultListOpts] to request the
// API's default page, then advance options with the returned list's
// [List.NextOffset] method:
//
//	opts := airplaneslive.DefaultListOpts
//	for {
//	    page, err := api.Airlines(ctx, "", "", "", "", "", opts)
//	    if err != nil {
//	        return err
//	    }
//	    // Do something with page, then
//	    if len(page.Results) == 0 {
//	        break
//	    }
//	    opts = page.NextOffset(opts)
//	}
//
// A basic aircraft query looks like:
//
//	ctx := context.Background()
//	api, err := airplaneslive.NewApi()
//	if err != nil {
//	    return err
//	}
//	point, err := geocoord.NewCoordinate(51.5074, -0.1278)
//	if err != nil {
//	    return err
//	}
//	response, err := api.AircraftWithinRadius(ctx, point, 10)
//	if err != nil {
//	    return err
//	}
//	_ = response.Ac
//
// The examples assume code is inside a function returning an error and imports
// context, github.com/fireflycons/airplaneslive, and github.com/fireflycons/geocoord.
// API responses may omit aircraft fields when data is unavailable.
package airplaneslive
