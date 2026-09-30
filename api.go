package airplaneslive

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"slices"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/fireflycons/geocoord"
)

var (
	normalHeaders = http.Header{
		"accept":             {"application/json"},
		"accept-encoding":    {"gzip, deflate"},
		"accept-language":    {"en-GB,en-US;q=0.9,en;q=0.8"},
		"origin":             {"https://airplanes.live"},
		"priority":           {"u=1, i"},
		"referer":            {"https://airplanes.live/"},
		"user-agent":         {"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/150.0.0.0 Safari/537.36"},
		"sec-ch-ua":          {`"Google Chrome";v="150", "Not_A Brand";v="8", "Chromium";v="150"`},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {`"Linux"`},
		"sec-fetch-dest":     {"empty"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-site":     {"same-site"},
		http.HeaderOrderKey: {
			"accept",
			"accept-encoding",
			"accept-language",
			"user-agent",
		},
	}
)

type errorResponse struct {
	Message string `json:"error"`
}

// Api provides access to the Airplanes.live aircraft and reference-data APIs.
type Api struct {
	client  tls_client.HttpClient
	baseUrl string
}

func newApi(baseURl string) (*Api, error) {
	options := []tls_client.HttpClientOption{
		tls_client.WithClientProfile(profiles.Chrome_150),
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithNotFollowRedirects(),
	}
	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	return &Api{
		client:  client,
		baseUrl: baseURl,
	}, nil
}

// NewApi creates an API client configured for Airplanes.live requests.
func NewApi() (*Api, error) {
	return newApi("https://api.airplanes.live")
}

// AircraftWithinRadius returns aircraft within radius nautical miles of coord.
func (api *Api) AircraftWithinRadius(ctx context.Context, coord geocoord.Coordinate, radius float64) (*V2Response, error) {

	if radius <= 0 {
		return nil, fmt.Errorf("radius must be greater than 0")
	}

	url := fmt.Sprintf("%s/v2/point/%f/%f/%f", api.baseUrl, coord.Latitude(), coord.Longitude(), radius)

	data, err := api.doRequest(ctx, url, normalHeaders)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch aircraft data: %w", err)
	}

	var response V2Response
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// AircraftByReg returns aircraft positions matching the supplied registrations.
func (api *Api) AircraftByReg(ctx context.Context, reg []string) (*V2Response, error) {

	param := strings.ReplaceAll(strings.Join(reg, ","), " ", "")

	url := fmt.Sprintf("%s/v2/reg/%s", api.baseUrl, param)

	data, err := api.doRequest(ctx, url, normalHeaders)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch aircraft data: %w", err)
	}

	var response V2Response
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// TaggedMilitaryAircraft returns all broadcasting miltary aircraft worldwide.
// Which fields are populated varies, likely depending on what each airforce is prepared to disclose.
func (api *Api) TaggedMilitaryAircraft(ctx context.Context) (*V2Response, error) {

	data, err := api.doRequest(ctx, api.baseUrl+"/v2/mil", normalHeaders)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch military aircraft data: %w", err)
	}

	var response V2Response
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// Airports searches or lists airports matching the supplied filters.
// An empty string for a filter repesents all possible values for that filter.
func (api *Api) Airports(ctx context.Context, cityCode, countryCode, iataCode, icaoCode string, opts ListOpts) (*AirportList, error) {

	qry := buildListQuery(
		map[string]string{
			"city_code":    cityCode,
			"country_code": countryCode,
			"iata_code":    iataCode,
			"icao_code":    icaoCode,
		},
		opts,
	)

	url := fmt.Sprintf("%s/rest/v1/ref/airports?%s", api.baseUrl, qry.Encode())

	var response AirportList

	if err := api.doListRequest(ctx, url, &response); err != nil {
		return nil, fmt.Errorf("failed to fetch airports data: %w", err)
	}

	return &response, nil
}

// Airlines searches or lists airlines matching the supplied filters.
// An empty string for a filter repesents all possible values for that filter.
func (api *Api) Airlines(ctx context.Context, name, callsign, countryCode, iataCode, icaoCode string, opts ListOpts) (*AirlineList, error) {

	qry := buildListQuery(
		map[string]string{
			"name":         name,
			"callsign":     callsign,
			"country_code": countryCode,
			"iata_code":    iataCode,
			"icao_code":    icaoCode,
		},
		opts,
	)

	url := fmt.Sprintf("%s/rest/v1/ref/airlines?%s", api.baseUrl, qry.Encode())

	var response AirlineList

	if err := api.doListRequest(ctx, url, &response); err != nil {
		return nil, fmt.Errorf("failed to fetch airlines data: %w", err)
	}

	return &response, nil
}

// Countries searches or lists countries matching the supplied filters.
// An empty string for a filter repesents all possible values for that filter.
func (api *Api) Countries(ctx context.Context, code, code3, continent string, opts ListOpts) (*CountryList, error) {

	qry := buildListQuery(
		map[string]string{
			"code":      code,
			"code3":     code3,
			"continent": continent,
		},
		opts,
	)

	url := fmt.Sprintf("%s/rest/v1/ref/countries?%s", api.baseUrl, qry.Encode())

	var response CountryList

	if err := api.doListRequest(ctx, url, &response); err != nil {
		return nil, fmt.Errorf("failed to fetch country data: %w", err)
	}

	return &response, nil
}

// Cities searches or lists cities matching the supplied filters.
// An empty string for a filter repesents all possible values for that filter.
func (api *Api) Cities(ctx context.Context, cityCode, countryCode string, opts ListOpts) (*CityList, error) {

	qry := buildListQuery(
		map[string]string{
			"city_code":    cityCode,
			"country_code": countryCode,
		},
		opts,
	)

	url := fmt.Sprintf("%s/rest/v1/ref/cities?%s", api.baseUrl, qry.Encode())

	var response CityList

	if err := api.doListRequest(ctx, url, &response); err != nil {
		return nil, fmt.Errorf("failed to fetch cities data: %w", err)
	}

	return &response, nil

}

// Timezones lists all timezone records.
func (api *Api) Timezones(ctx context.Context, opts ListOpts) (*TimezoneList, error) {

	qry := buildListQuery(nil, opts)
	url := fmt.Sprintf("%s/rest/v1/ref/timezones?%s", api.baseUrl, qry.Encode())

	var response TimezoneList
	if err := api.doListRequest(ctx, url, &response); err != nil {
		return nil, fmt.Errorf("failed to fetch timezone data: %w", err)
	}

	return &response, nil
}

func (api *Api) doRequest(ctx context.Context, url string, headers http.Header) ([]byte, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header = headers

	resp, err := api.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if slices.Contains([]int{400, 502, 503}, resp.StatusCode) {
		var errResp errorResponse
		body, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &errResp); err != nil {
			return nil, fmt.Errorf("request failed with status code %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, errResp.Message)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

func (api *Api) doListRequest(ctx context.Context, uri string, list any) error {

	data, err := api.doRequest(ctx, uri, normalHeaders)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, list); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

func buildListQuery(params map[string]string, opts ListOpts) url.Values {

	qry := url.Values{}

	for k, v := range params {
		if v != "" {
			qry.Set(k, v)
		}
	}

	if opts.offset != 0 {
		qry.Set("offset", fmt.Sprintf("%d", opts.offset))
	}
	if opts.limit != 0 {
		qry.Set("limit", fmt.Sprintf("%d", opts.limit))
	}

	fld := strings.ReplaceAll(strings.Join(opts.fields, ","), " ", "")
	if fld != "" {
		qry.Set("_fields", fld)
	}

	return qry
}
