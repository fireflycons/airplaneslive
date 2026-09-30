package airplaneslive

// List contains a page of results of type T and the API's pagination metadata.
type List[T any] struct {
	// Total number of items in the list.
	Count int `json:"count"`

	// Value of limit passed to API call, or default limit if not specified.
	Limit int `json:"limit"`

	// Value of offset passed to API call, or default offset if not specified.
	Offset int `json:"offset"`

	// Array of results of type T.
	Results []T `json:"results"`
}

// NextOffset returns options advanced by the number of results in the page.
func (l *List[T]) NextOffset(o ListOpts) ListOpts {

	o.offset += len(l.Results)
	return o
}
