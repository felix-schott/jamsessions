package geocoding

import (
	"log"
	"math"
	"testing"

	geom "github.com/twpayne/go-geom"
)

type address struct {
	Street   string
	City     string
	Postcode string
}

func TestGeocoding(t *testing.T) {
	cases := map[*address]*geom.Point{
		{Street: "47 Frith Street", City: "London", Postcode: "W1D 4HT"}: geom.NewPoint(geom.XY).MustSetCoords([]float64{-0.132, 51.513}),
		{Street: "6 Moor Street", City: "London", Postcode: "W1D 5NA"}:   geom.NewPoint(geom.XY).MustSetCoords([]float64{-0.12980699914457172, 51.51339645}),
	}

	for tc, exp := range cases {
		result, err := Geocode(tc.Street, tc.City, tc.Postcode)
		if err != nil {
			err, ok := err.(NominatimDownError)
			if ok {
				log.Println("nominatim unavailable, skipping test:", err)
				return
			} else {
				t.Fatalf("an error occured when geocoding: %v", err)
			}
		}
		if math.Round(result.Y()*1000)/1000 != math.Round(exp.Y()*1000)/1000 {
			t.Errorf("property Lat is different (exp: %v, got: %v)", exp.Y(), result.Y())
		}
		if math.Round(result.X()*1000)/1000 != math.Round(exp.X()*1000)/1000 {
			t.Errorf("property Lon is different (exp: %v, got: %v)", exp.X(), result.X())
		}
	}
}
