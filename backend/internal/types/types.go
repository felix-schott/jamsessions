package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type ValidationError struct {
	Msg string
}

func (e ValidationError) Error() string {
	return e.Msg
}

// GENRE ENUM
type Genre string

func (g Genre) String() string {
	return strings.Trim(string(g), `"`)
}

// values must match database schema constraint
const (
	Any           Genre = "Any"
	StraightAhead Genre = "Straight-Ahead_Jazz"
	JazzFunk      Genre = "Jazz-Funk"
	Fusion        Genre = "Fusion"
	LatinJazz     Genre = "Latin_Jazz"
	ModernJazz    Genre = "Modern_Jazz"
	TradJazz      Genre = "Trad_Jazz"
	Funk          Genre = "Funk"
	RnB           Genre = "RnB"
	HipHop        Genre = "Hip-Hop"
	Blues         Genre = "Blues"
	Folk          Genre = "Folk"
	Rock          Genre = "Rock"
	Pop           Genre = "Pop"
	WorldMusic    Genre = "World_Music"
)

var Genres = map[Genre]struct{}{ // an empty struct doesn't occupy any bytes in memory, good way to emulate a set
	Any:           {},
	StraightAhead: {},
	JazzFunk:      {},
	Fusion:        {},
	LatinJazz:     {},
	ModernJazz:    {},
	TradJazz:      {},
	Funk:          {},
	RnB:           {},
	HipHop:        {},
	Blues:         {},
	Folk:          {},
	Rock:          {},
	Pop:           {},
	WorldMusic:    {},
}

func (g *Genre) UnmarshalJSON(b []byte) error {
	s := Genre(strings.Trim(string(b), `"`))
	if _, ok := Genres[s]; !ok {
		var validGenres = make([]string, len(Genres))
		i := 0
		for k := range Genres {
			validGenres[i] = k.String()
			i += 1
		}
		return ValidationError{Msg: fmt.Sprintf("%s is not a valid Genre. Valid values: %v", b, strings.Join(validGenres, ", "))}
	}
	*g = s
	return nil
}

func (g Genre) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.String())
}

// BACKLINE ENUM

type Backline string

func (b Backline) String() string {
	return strings.Trim(string(b), `"`)
}

const (
	PA             Backline = "PA"
	GuitarAmp      Backline = "Guitar_Amp"
	BassAmp        Backline = "Bass_Amp"
	Drums          Backline = "Drums"
	Keys           Backline = "Keys"
	Mic            Backline = "Microphone"
	MiscPercussion Backline = "MiscPercussion"
)

var BacklineOptions = map[Backline]struct{}{ // an empty struct doesn't occupy any bytes in memory, good way to emulate a set
	PA:             {},
	GuitarAmp:      {},
	BassAmp:        {},
	Drums:          {},
	Keys:           {},
	Mic:            {},
	MiscPercussion: {},
}

func (backline *Backline) UnmarshalJSON(b []byte) error {
	s := Backline(strings.Trim(string(b), `"`))
	if _, ok := BacklineOptions[s]; !ok {
		var validBacklineOptions = make([]string, len(BacklineOptions))
		i := 0
		for k := range BacklineOptions {
			validBacklineOptions[i] = k.String()
			i += 1
		}
		return ValidationError{Msg: fmt.Sprintf("%s is not a valid Backline option. Valid values: %v", b, strings.Join(validBacklineOptions, ", "))}
	}
	*backline = s
	return nil
}

func (backline Backline) MarshalJSON() ([]byte, error) {
	return json.Marshal(backline.String())
}

// INTERVAL ENUM
type Interval string

func (i Interval) String() string {
	return strings.Trim(string(i), `"`)
}

const (
	Once            Interval = "Once"
	Daily           Interval = "Daily"
	Weekly          Interval = "Weekly"
	Fortnightly     Interval = "Fortnightly"
	FirstOfMonth    Interval = "FirstOfMonth"
	SecondOfMonth   Interval = "SecondOfMonth"
	ThirdOfMonth    Interval = "ThirdOfMonth"
	FourthOfMonth   Interval = "FourthOfMonth"
	LastOfMonth     Interval = "LastOfMonth"
	IrregularWeekly Interval = "IrregularWeekly"
)

var IntervalOptions = map[Interval]struct{}{ // an empty struct doesn't occupy any bytes in memory, good way to emulate a set
	Once:            {},
	Daily:           {},
	Weekly:          {},
	Fortnightly:     {},
	FirstOfMonth:    {},
	SecondOfMonth:   {},
	ThirdOfMonth:    {},
	FourthOfMonth:   {},
	LastOfMonth:     {},
	IrregularWeekly: {},
}

func (i *Interval) UnmarshalJSON(b []byte) error {
	s := Interval(strings.Trim(string(b), `"`))
	if _, ok := IntervalOptions[s]; !ok {
		var validIntervalOptions = make([]string, len(IntervalOptions))
		i := 0
		for k := range IntervalOptions {
			validIntervalOptions[i] = k.String()
			i += 1
		}
		return ValidationError{Msg: fmt.Sprintf("%s is not a valid Interval option. Valid values: %v", b, strings.Join(validIntervalOptions, ", "))}
	}
	*i = s
	return nil
}

func (i Interval) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// TYPE DATE
type Date time.Time

func (d *Date) UnmarshalJSON(b []byte) error {
	t, err := time.Parse("2006-01-02", strings.Trim(string(b), `"`))
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d).Format(time.DateOnly))
}

func (d Date) Format(s string) string {
	t := time.Time(d)
	return t.Format(s)
}

func (d Date) String() string {
	return d.Format(time.DateOnly)
}

// GEOJSON

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type VenueProperties struct {
	VenueID           *int32      `json:"venue_id,omitempty"`
	VenueName         *string     `json:"venue_name,omitempty"`
	AddressFirstLine  *string     `json:"address_first_line,omitempty"`
	AddressSecondLine *string     `json:"address_second_line,omitempty"`
	City              *string     `json:"city,omitempty"`
	Postcode          *string     `json:"postcode,omitempty"`
	VenueWebsite      *string     `json:"venue_website,omitempty"`
	Backline          *[]Backline `json:"backline,omitempty"`
	VenueComments     *[]string   `json:"venue_comments,omitempty"`
	VenueDtUpdatedUtc *time.Time  `json:"venue_dt_updated_utc,omitempty"`
}

type VenueFeature struct {
	Type       string          `json:"type"`
	Properties VenueProperties `json:"properties"`
	Geometry   Geometry        `json:"geometry"`
}

type SessionProperties struct {
	SessionID       *int32     `json:"session_id,omitempty"`
	SessionName     *string    `json:"session_name,omitempty"`
	Venue           *int32     `json:"venue,omitempty"`
	Description     *string    `json:"description,omitempty"`
	Genres          *[]Genre   `json:"genres,omitempty"`
	StartTimeUtc    *time.Time `json:"start_time_utc,omitempty"`
	Interval        *Interval  `json:"interval,omitempty"`
	DurationMinutes *int16     `json:"duration_minutes,omitempty"`
	SessionWebsite  *string    `json:"session_website,omitempty"`
	DtUpdatedUtc    *time.Time `json:"dt_updated_utc,omitempty"`
	Rating          *float32   `json:"rating,omitempty"`
	Dates           *[]Date    `json:"dates,omitempty"`
}

type SessionPropertiesWithVenue struct {
	SessionProperties
	VenueProperties
}

type SessionFeature[T SessionProperties | SessionPropertiesWithVenue] struct {
	Type       string   `json:"type"`
	Properties T        `json:"properties"`
	Geometry   Geometry `json:"geometry"`
}

type FeatureCollection[T SessionFeature[SessionProperties] | SessionFeature[SessionPropertiesWithVenue] | VenueFeature] struct {
	Type     string `json:"type"`
	Features []T    `json:"features"`
}

// declare alias
type SessionFeatureCollection = FeatureCollection[SessionFeature[SessionProperties]]
type SessionWithVenueFeatureCollection = FeatureCollection[SessionFeature[SessionPropertiesWithVenue]]
type VenueFeatureCollection = FeatureCollection[VenueFeature]
