package airplaneslive

// Country represents a country entity as defined by the Airplanes.live API.
type Country struct {
	// ISO 3166-1 alpha-2 code (required).
	Code string `json:"code"`

	// ISO 3166-1 alpha-3 code.
	Code3 *string `json:"code3,omitempty"`

	// Full name of the country (required).
	Name string `json:"name"`

	// Latest available population figure.
	Population *int64 `json:"population,omitempty"`

	// Continent name.
	Continent *string `json:"continent,omitempty"`

	// ISO 4217 currency code.
	Currency *string `json:"currency,omitempty"`
}

// CountryList contains a page of country records.
type CountryList List[Country]

// NextOffset returns list options advanced by the number of results in the page.
func (l *CountryList) NextOffset(o ListOpts) ListOpts {
	return (*List[Country])(l).NextOffset(o)
}
