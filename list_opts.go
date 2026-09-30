package airplaneslive

// DefaultListOpts requests the API defaults and all available fields.
// The API generally uses offset zero and a limit of 50 when these are omitted.
var DefaultListOpts = ListOpts{}

// ListOpts specifies pagination and field selection for reference-data requests.
type ListOpts struct {
	offset int
	limit  int
	fields []string
}

// NewListOpts creates list options with the given limit, offset, and selected fields.
func NewListOpts(limit, offset int, fields ...string) ListOpts {
	return ListOpts{
		offset: offset,
		limit:  limit,
		fields: fields,
	}
}
