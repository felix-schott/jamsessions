package geocoding

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func serviceIsHealthy() bool {
	resp, err := client.Get("https://nominatim.openstreetmap.org/status")
	if err != nil {
		return false
	}
	return resp.StatusCode == 200
}

func Geocode(street string, city string, postcode string) (*geom.Point, error) {
	if !serviceIsHealthy() {
		return nil, errors.New("nominatim API is down")
	}

	reqUrl := fmt.Sprintf("https://nominatim.openstreetmap.org/search?street=%v&city=%v&country=UK&postcode=%v&format=geojson&limit=1", url.QueryEscape(street), url.QueryEscape(city), url.QueryEscape(postcode))
	resp, err := client.Get(reqUrl)
	if err != nil {
		return nil, fmt.Errorf("an unkown error occured when making request to %v: %v", reqUrl, err)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading the body: %v", err)
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
