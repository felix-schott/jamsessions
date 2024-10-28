package types

import (
	"time"
)

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type VenueProperties struct {
	VenueID           *int32     `json:"venue_id,omitempty"`
	VenueName         *string    `json:"venue_name,omitempty"`
	AddressFirstLine  *string    `json:"address_first_line,omitempty"`
	AddressSecondLine *string    `json:"address_second_line,omitempty"`
	City              *string    `json:"city,omitempty"`
	Postcode          *string    `json:"postcode,omitempty"`
	VenueWebsite      *string    `json:"venue_website,omitempty"`
	Backline          *[]string  `json:"backline,omitempty"`
	VenueComments     *[]string  `json:"venue_comments,omitempty"`
	VenueDtUpdatedUtc *time.Time `json:"venue_dt_updated_utc,omitempty"`
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
	Genres          *[]string  `json:"genres,omitempty"`
	StartTimeUtc    *time.Time `json:"start_time_utc,omitempty"`
	Interval        *string    `json:"interval,omitempty"`
	DurationMinutes *int16     `json:"duration_minutes,omitempty"`
	SessionWebsite  *string    `json:"session_website,omitempty"`
	DtUpdatedUtc    *time.Time `json:"dt_updated_utc,omitempty"`
	Rating          *float32   `json:"rating,omitempty"`
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
