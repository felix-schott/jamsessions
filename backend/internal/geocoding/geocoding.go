package geocoding

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	geom "github.com/twpayne/go-geom"
	"golang.org/x/time/rate"
)

var client *httpClientWithRateLimit

type geometry struct {
	Type        string
	Coordinates []float64
}

type feature struct {
	// we don't care about the rest of the fields
	Geometry geometry `json:"geometry"`
}

type nominatimResponse struct {
	Type     string
	Licence  string
	Features []feature
}

type NominatimDownError struct {
	StatusCode int
	Body       []byte
	Err        error
}

func (r NominatimDownError) Error() string {
	return fmt.Sprintf("Nominatim unavailable (http status: %d, body: %s, error: %v)", r.StatusCode, r.Body, r.Err)
}

// Returns a NominatimDownError if the service is not healthy, otherwise nil
func serviceIsHealthy() error {
	req, err := http.NewRequest("GET", "https://nominatim.openstreetmap.org/status", nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	var b []byte
	if err != nil {
		b, err = io.ReadAll(resp.Body)
		if err != nil {
			return NominatimDownError{StatusCode: resp.StatusCode, Err: err}
		}
		return NominatimDownError{StatusCode: resp.StatusCode, Body: b}
	}
	if resp.StatusCode != 200 {
		b, err = io.ReadAll(resp.Body)
		if err != nil {
			return NominatimDownError{StatusCode: resp.StatusCode, Err: err}
		}
		return NominatimDownError{StatusCode: resp.StatusCode, Body: b}
	}
	return nil
}

func Geocode(street string, city string, postcode string) (*geom.Point, error) {
	// instantiate client that respects nominatim rate limit (max 1 request per second)
	if client == nil {
		rl := rate.NewLimiter(rate.Every(time.Second*1), 1)
		client = NewHttpClient(rl, "github.com/felix-schott/jamsessions")
	}

	err := serviceIsHealthy()
	if err != nil {
		return nil, err
	}

	reqUrl := fmt.Sprintf("https://nominatim.openstreetmap.org/search?street=%v&city=%v&country=UK&postcode=%v&format=geojson&limit=1", url.QueryEscape(street), url.QueryEscape(city), url.QueryEscape(postcode))
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}
	resp, err2 := client.Do(req)
	if err2 != nil {
		return nil, fmt.Errorf("an unkown error occured when making request to %v: %v", reqUrl, err2)
	}
	bodyBytes, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		return nil, fmt.Errorf("error reading the body: %v", err2)
	}
	if resp.StatusCode != 200 {
		body := string(bodyBytes)
		return nil, fmt.Errorf("encountered a non-200 status code when making request to %v: %v (Body: %v)", reqUrl, resp.Status, body)
	}
	var result nominatimResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil { // Parse []byte to result pointer
		return nil, fmt.Errorf("failed to unmarshal Nominatim response: %v", err)
	}
	if len(result.Features) == 0 {
		return nil, fmt.Errorf("no matches for %v, %v, %v (url %v)", street, city, postcode, reqUrl)
	}

	return geom.NewPoint(geom.XY).MustSetCoords(result.Features[0].Geometry.Coordinates).SetSRID(4326), nil
}
