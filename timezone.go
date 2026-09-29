package airplaneslive

// Timezone represents a timezone entity as defined by the Airplanes.live API.
type Timezone struct {
	// IANA timezone id (required).
	Timezone string `json:"timezone"`

	// ISO 3166-1 alpha-2 code.
	CountryCode *string `json:"country_code,omitempty"`

	// GMT offset in hours (standard time).
	GMT *float64 `json:"gmt,omitempty"`

	// Offset in hours during DST.
	DST *float64 `json:"dst,omitempty"`
}

// TimezoneList contains a page of timezone records.
type TimezoneList List[Timezone]

// NextOffset returns list options advanced by the number of results in the page.
func (l *TimezoneList) NextOffset(o ListOpts) ListOpts {
	return (*List[Timezone])(l).NextOffset(o)
}
