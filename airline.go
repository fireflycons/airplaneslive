package airplaneslive

// Airline represents airline identification and reference details.
type Airline struct {
	// Name of the airline.
	Name string `json:"name"`

	// IATA code (2 characters, optional/nullable).
	IATACode *string `json:"iata_code,omitempty"`

	// ICAO code (3 characters, optional/nullable).
	ICAOCode *string `json:"icao_code,omitempty"`

	// Radio callsign (optional/nullable).
	Callsign *string `json:"callsign,omitempty"`

	// ISO 3166-1 alpha-2 country code (optional/nullable).
	CountryCode *string `json:"country_code,omitempty"`

	// Official website URL (optional/nullable).
	Website *string `json:"website,omitempty"`
}

// AirlineList contains a page of airline records.
type AirlineList List[Airline]

// NextOffset returns list options advanced by the number of results in the page.
func (l *AirlineList) NextOffset(o ListOpts) ListOpts {
	return (*List[Airline])(l).NextOffset(o)
}
