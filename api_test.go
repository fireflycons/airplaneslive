package airplaneslive

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/fireflycons/geocoord"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApiEndpoints(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		query          url.Values
		call           func(*Api) (any, error)
		responseBody   string
		assertResponse func(*testing.T, any)
	}{
		{
			name:  "aircraft within radius",
			path:  "/v2/point/51.500000/-0.120000/25.000000",
			query: url.Values{},
			call: func(api *Api) (any, error) {
				return api.AircraftWithinRadius(t.Context(), geocoord.MustNewCoordinate(51.5, -0.12), 25)
			},
			responseBody: `{"ac":[],"msg":"ready","now":0,"total":2}`,
			assertResponse: func(t *testing.T, response any) {
				result, ok := response.(*V2Response)
				require.True(t, ok)
				assert.Equal(t, "ready", result.Msg)
				assert.Equal(t, 2, result.Total)
			},
		},
		{
			name:  "aircraft by registrations",
			path:  "/v2/reg/N123AB,G-ABCD",
			query: url.Values{},
			call: func(api *Api) (any, error) {
				return api.AircraftByReg(t.Context(), []string{"N123AB", " G-ABCD "})
			},
			responseBody: `{"ac":[],"msg":"ready","now":0,"total":2}`,
			assertResponse: func(t *testing.T, response any) {
				result, ok := response.(*V2Response)
				require.True(t, ok)
				assert.Equal(t, "ready", result.Msg)
				assert.Equal(t, 2, result.Total)
			},
		},
		{
			name:  "tagged military aircraft",
			path:  "/v2/mil",
			query: url.Values{},
			call: func(api *Api) (any, error) {
				return api.TaggedMilitaryAircraft(t.Context())
			},
			responseBody: `{"ac":[],"msg":"ready","now":0,"total":2}`,
			assertResponse: func(t *testing.T, response any) {
				result, ok := response.(*V2Response)
				require.True(t, ok)
				assert.Equal(t, "ready", result.Msg)
				assert.Equal(t, 2, result.Total)
			},
		},
		{
			name:  "airports",
			path:  "/rest/v1/ref/airports",
			query: url.Values{"city_code": {"LON"}, "country_code": {"GB"}, "iata_code": {"LHR"}, "icao_code": {"EGLL"}, "limit": {"25"}, "offset": {"50"}, "_fields": {"name,country_code"}},
			call: func(api *Api) (any, error) {
				return api.Airports(t.Context(), "LON", "GB", "LHR", "EGLL", NewListOpts(25, 50, "name", "country_code"))
			},
			responseBody: `{"count":7,"limit":25,"offset":50,"results":[{"name":"Test Airport"}]}`,
			assertResponse: func(t *testing.T, response any) {
				result, ok := response.(*AirportList)
				require.True(t, ok)
				assert.Equal(t, 7, result.Count)
				require.Len(t, result.Results, 1)
				assert.Equal(t, "Test Airport", result.Results[0].Name)
			},
		},
		{
			name:  "airlines",
			path:  "/rest/v1/ref/airlines",
			query: url.Values{"name": {"Air Example"}, "callsign": {"EXAMPLE"}, "country_code": {"GB"}, "iata_code": {"EX"}, "icao_code": {"EXA"}, "limit": {"25"}, "offset": {"50"}, "_fields": {"name,country_code"}},
			call: func(api *Api) (any, error) {
				return api.Airlines(t.Context(), "Air Example", "EXAMPLE", "GB", "EX", "EXA", NewListOpts(25, 50, "name", "country_code"))
			},
			responseBody: `{"count":7,"limit":25,"offset":50,"results":[]}`,
			assertResponse: func(t *testing.T, response any) {
				result, ok := response.(*AirlineList)
				require.True(t, ok)
				assert.Equal(t, 7, result.Count)
			},
		},
		{
			name:  "countries",
			path:  "/rest/v1/ref/countries",
			query: url.Values{"code": {"GB"}, "code3": {"GBR"}, "continent": {"EU"}, "limit": {"25"}, "offset": {"50"}, "_fields": {"name,country_code"}},
			call: func(api *Api) (any, error) {
				return api.Countries(t.Context(), "GB", "GBR", "EU", NewListOpts(25, 50, "name", "country_code"))
			},
			responseBody: `{"count":7,"limit":25,"offset":50,"results":[]}`,
			assertResponse: func(t *testing.T, response any) {
				result, ok := response.(*CountryList)
				require.True(t, ok)
				assert.Equal(t, 7, result.Count)
			},
		},
		{
			name:  "cities",
			path:  "/rest/v1/ref/cities",
			query: url.Values{"city_code": {"LON"}, "country_code": {"GB"}, "limit": {"25"}, "offset": {"50"}, "_fields": {"name,country_code"}},
			call: func(api *Api) (any, error) {
				return api.Cities(t.Context(), "LON", "GB", NewListOpts(25, 50, "name", "country_code"))
			},
			responseBody: `{"count":7,"limit":25,"offset":50,"results":[]}`,
			assertResponse: func(t *testing.T, response any) {
				result, ok := response.(*CityList)
				require.True(t, ok)
				assert.Equal(t, 7, result.Count)
			},
		},
		{
			name:  "timezones",
			path:  "/rest/v1/ref/timezones",
			query: url.Values{"limit": {"25"}, "offset": {"50"}, "_fields": {"name,country_code"}},
			call: func(api *Api) (any, error) {
				return api.Timezones(t.Context(), NewListOpts(25, 50, "name", "country_code"))
			},
			responseBody: `{"count":7,"limit":25,"offset":50,"results":[]}`,
			assertResponse: func(t *testing.T, response any) {
				result, ok := response.(*TimezoneList)
				require.True(t, ok)
				assert.Equal(t, 7, result.Count)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type request struct {
				path  string
				query url.Values
			}
			requests := make(chan request, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests <- request{path: r.URL.Path, query: r.URL.Query()}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			api, err := newApi(server.URL)
			require.NoError(t, err)

			response, err := tt.call(api)
			require.NoError(t, err)
			tt.assertResponse(t, response)

			gotRequest := <-requests
			assert.Equal(t, tt.path, gotRequest.path)
			assert.Equal(t, tt.query, gotRequest.query)
		})
	}
}

func TestDoRequestErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "API error body", statusCode: 400, body: `{"error":"bad request"}`, wantError: "request failed with status code 400: bad request"},
		{name: "invalid API error body", statusCode: 502, body: "not json", wantError: "request failed with status code 502"},
		{name: "other client error", statusCode: 401, body: `{"error":"unauthorized"}`, wantError: "request failed with status code 401"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()

			api, err := newApi(server.URL)
			require.NoError(t, err)

			_, err = api.doRequest(t.Context(), server.URL, nil)
			require.EqualError(t, err, tt.wantError)
		})
	}
}

func TestAircraftWithinRadiusRejectsNonPositiveRadius(t *testing.T) {
	api := &Api{}
	for _, radius := range []float64{0, -1} {
		_, err := api.AircraftWithinRadius(t.Context(), geocoord.MustNewCoordinate(0, 0), radius)
		assert.EqualError(t, err, "radius must be greater than 0")
	}
}
