package geocoding

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	geom "github.com/twpayne/go-geom"
)

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

var tr = &http.Transport{
	IdleConnTimeout: 30 * time.Second,
}
var client = &http.Client{Transport: tr}

// nominatim
// 1 request per second!!

type NominatimDownError struct {
	StatusCode int
	Body       []byte
}

func (r NominatimDownError) Error() string {
	return fmt.Sprintf("Nominatim unavailable (http status: %d, body: %s)", r.StatusCode, r.Body)
}

func serviceIsHealthy() (bool, NominatimDownError) {
	resp, err := client.Get("https://nominatim.openstreetmap.org/status")
	if err != nil {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return false, NominatimDownError{StatusCode: resp.StatusCode, Body: []byte{}}
		}
		return false, NominatimDownError{StatusCode: resp.StatusCode, Body: b}
	}
	slog.Error("nominatimDown", "status", resp.StatusCode)
	return resp.StatusCode == 200, NominatimDownError{StatusCode: resp.StatusCode, Body: []byte{}}
}

func Geocode(street string, city string, postcode string) (*geom.Point, error) {
	healthy, err := serviceIsHealthy()
	if !healthy {
		return nil, err
	}

	reqUrl := fmt.Sprintf("https://nominatim.openstreetmap.org/search?street=%v&city=%v&country=UK&postcode=%v&format=geojson&limit=1", url.QueryEscape(street), url.QueryEscape(city), url.QueryEscape(postcode))
	resp, err2 := client.Get(reqUrl)
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
