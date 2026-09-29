package airplaneslive

// City represents a city entity as defined by the Airplanes.live API.
type City struct {
	// Name of the city (required).
	Name string `json:"name"`

	// IATA city code (required).
	CityCode string `json:"city_code"`

	// UN/LOCODE.
	UNLocode *string `json:"un_locode,omitempty"`

	// Latitude coordinate.
	Lat *float64 `json:"lat,omitempty"`

	// Longitude coordinate.
	Lon *float64 `json:"lon,omitempty"`

	// Elevation in metres.
	Alt *int `json:"alt,omitempty"`

	// IANA timezone.
	Timezone *string `json:"timezone,omitempty"`

	// ISO 3166-1 alpha-2 country code.
	CountryCode *string `json:"country_code,omitempty"`

	// Latest available population figure.
	Population *int64 `json:"population,omitempty"`

	// Wikipedia article URL.
	Wikipedia *string `json:"wikipedia,omitempty"`
}

// CityList contains a page of city records.
type CityList List[City]

// NextOffset returns list options advanced by the number of results in the page.
func (l *CityList) NextOffset(o ListOpts) ListOpts {
	return (*List[City])(l).NextOffset(o)
}
