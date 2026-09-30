# airplaneslive

`airplaneslive` provides a Go client for aircraft position and reference-data queries against the Airplanes.live API.

Create an `Api` with `NewApi`, then pass a context to each request. Aircraft queries include positions within a radius and positions by registration. Reference-data methods search or list airlines, airports, cities, countries, and timezones.

## Coordinates and Units

Coordinates are `geocoord.Coordinate` values from [github.com/fireflycons/geocoord](https://github.com/fireflycons/geocoord). Use `geocoord.NewCoordinate` to validate latitude in `[-90, 90]` and longitude in `[-180, 180]`, or `geocoord.MustNewCoordinate` to panic on invalid values. Aircraft distances and radius values are in nautical miles, and headings are in degrees.

## Aircraft Query

```go
ctx := context.Background()
api, err := airplaneslive.NewApi()
if err != nil {
    return err
}
point, err := geocoord.NewCoordinate(51.5074, -0.1278)
if err != nil {
    return err
}
response, err := api.AircraftWithinRadius(ctx, point, 10)
if err != nil {
    return err
}
_ = response.Ac
```

This example assumes it is inside a function returning an error. Import `context`, `github.com/fireflycons/airplaneslive`, and `github.com/fireflycons/geocoord` in your application.

## Reference-Data Pagination

Reference-data results are paginated. Pass `DefaultListOpts` to request the API's default page, then advance options with the returned list's `NextOffset` method:

```go
opts := airplaneslive.DefaultListOpts
for {
    page, err := api.Airlines(ctx, "", "", "", "", "", opts)
    if err != nil {
        return err
    }
    // Do something with page, then
    if len(page.Results) == 0 {
        break
    }
    opts = page.NextOffset(opts)
}
```

API responses may omit aircraft fields when data is unavailable.

## Package Documentation

See [here](./docs/package.md)
