package airplaneslive

// Airport represents an airport entity as defined by the Airplanes.live API.
type Airport struct {
	// Name of the airport (required).
	Name string `json:"name"`

	// IATA code (e.g., "JFK").
	IATACode *string `json:"iata_code"`

	// ICAO code (e.g., "KJFK").
	ICAOCode *string `json:"icao_code"`

	// Latitude coordinate.
	Lat *float64 `json:"lat"`

	// Longitude coordinate.
	Lon *float64 `json:"lon"`

	// Elevation in metres.
	Alt *int `json:"alt"`

	// City name.
	City *string `json:"city"`

	// IATA city code (from OpenTravelData).
	CityCode *string `json:"city_code"`

	// UN/LOCODE (from OpenTravelData).
	UNLocode *string `json:"un_locode"`

	// IANA timezone (from OpenTravelData).
	Timezone *string `json:"timezone"`

	// ISO 3166-1 alpha-2 country code.
	CountryCode *string `json:"country_code"`

	// Official website URL.
	Website *string `json:"website"`

	// Number of non-closed runways.
	Runways *int `json:"runways"`
}

// AirportList contains a page of airport records.
type AirportList List[Airport]

// NextOffset returns list options advanced by the number of results in the page.
func (l *AirportList) NextOffset(o ListOpts) ListOpts {
	return (*List[Airport])(l).NextOffset(o)
}
