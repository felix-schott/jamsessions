package dbutils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/felix-schott/jamsessions/backend/internal/types"
	"github.com/jackc/pgx/v5/pgtype"
	geom "github.com/twpayne/go-geom"
)

var queries *Queries
var ctx = context.Background()
var fixtureWeeklySunday int32
var fixtureDaily int32
var fixtureThirdOfMonthMonday int32
var fixtureLastOfMonthSaturday int32
var fixtureFortnightlyWednesday int32

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	// teardown() // we use an ephemeral testing database (see Makefile) - no need for teardown
	os.Exit(code)
}

func setup() {
	pool, err := CreatePool(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// defer pool.Close()

	queries = New(pool)

	var numTables int
	row := pool.QueryRow(ctx, "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'london_jam_sessions' AND table_name IN ('jamsessions', 'venues');")
	if err != nil {
		log.Fatal(err)
	}
	row.Scan(&numTables)
	if numTables != 2 {
		log.Fatalf("test expects a pre-populated DB. Tables 'jamsessions' and 'venues' found missing. Number of tables: %v", numTables)
	}

	var numRowsPostgisFunc int
	row = pool.QueryRow(ctx, "SELECT COUNT(*) FROM information_schema.routines WHERE lower(routine_name::text) = 'st_asgeojson';")
	if err != nil {
		log.Fatal(err)
	}
	row.Scan(&numRowsPostgisFunc)
	if numRowsPostgisFunc == 0 {
		log.Fatal("test expects a postgis functions to be available. can't use 'st_asgeojson'.")
	}

	queries = New(pool)

	// add dummy records to both Venue and JamSession table
	fixtureWeeklySunday, err = insertSession("Ronnie Scott's Jazz Club", 12, 18, 15, "Weekly") // sunday
	if err != nil {
		log.Fatal(err)
	}

	fixtureDaily, err = insertSession("Ronnie Scott's Jazz Cafe", 13, 19, 15, "Daily")
	if err != nil {
		log.Fatal(err)
	}

	fixtureThirdOfMonthMonday, err = insertSession("Ronnie Scott's Jazz Bar", 14, 19, 15, "ThirdOfMonth") // monday
	if err != nil {
		log.Fatal(err)
	}

	fixtureLastOfMonthSaturday, err = insertSession("Ronnie Scott's Jazz Restaurant", 15, 24, 15, "LastOfMonth") // saturday
	if err != nil {
		log.Fatal(err)
	}

	fixtureFortnightlyWednesday, err = insertSession("Ronnie Scott's Jazz Diner", 16, 28, 15, "Fortnightly") // wednesday
	if err != nil {
		log.Fatal(err)
	}
}

// don't really need this as we're using an ephemeral testing database
// func teardown() {
// 	if fixtureWeeklySunday != 0 {
// 		err := queries.DeleteVenueByJamSessionId(ctx, fixtureWeeklySunday)
// 		if err != nil {
// 			log.Fatalf("could not delete venue previously inserted for session id %v: %v", fixtureWeeklySunday, err)
// 		}
// 		err = queries.DeleteJamSessionById(ctx, fixtureWeeklySunday)
// 		if err != nil {
// 			log.Fatalf("could not delete session with id %v", fixtureWeeklySunday)
// 		}
// 	}
// }

// https://github.com/golang/go/issues/63309
func ptr[T any](t T) *T { return &t }

// Helper function that inserts a JamSession and corresponding Location record.
// The two params houseNumber and startMinute can be used to control uniqueness within the
// Location and JamSession table, respectively.
func insertSession(locationName string, houseNumber uint8, day int, startMinute int, interval string) (int32, error) {
	venue, err := queries.InsertVenue(ctx, InsertVenueParams{
		VenueName:        locationName,
		AddressFirstLine: string(houseNumber) + " Frith Street",
		City:             "London",
		Postcode:         "W1D 4HT",
		Geom:             geom.NewPoint(geom.XY).MustSetCoords([]float64{-0.132, 51.513}),
		VenueWebsite:     ptr("https://www.ronniescotts.co.uk/"),
		Backline:         []string{"PA", "Drums"},
		VenueComments:    []string{"Sign up with host at the start of the session", "Very touristy"},
	})
	if err != nil {
		return 0, fmt.Errorf("the following error occurred when trying to insert a Location record: %v", err)
	}

	sessionId, err := queries.InsertJamSession(ctx, InsertJamSessionParams{
		SessionName:     "test_session",
		StartTimeUtc:    pgtype.Timestamptz{Time: time.Date(2024, 8, day, 19, startMinute, 0, 0, time.UTC), Valid: true},
		Interval:        interval,
		DurationMinutes: 120,
		Venue:           venue,
	})
	if err != nil {
		return 0, fmt.Errorf("the following error occured when trying to insert a Jamsession record: %v", err)
	}

	// insert two ratings
	_, err = queries.InsertSessionRating(ctx, InsertSessionRatingParams{
		Session: sessionId,
		Rating:  ptr(int16(5)),
	})
	if err != nil {
		return 0, fmt.Errorf("the following error occurred when trying to insert a rating: %v", err)
	}

	_, err = queries.InsertSessionRating(ctx, InsertSessionRatingParams{
		Session: sessionId,
		Rating:  ptr(int16(1)),
	})
	if err != nil {
		return 0, fmt.Errorf("the following error occurred when trying to insert a rating: %v", err)
	}

	return sessionId, nil
}

func TestInsertJamsession(t *testing.T) {
	var houseNumber uint8 = 50
	var startMinute int = 30
	newSessionId, err := insertSession("Ronnie Scott's Blues Club", houseNumber, 1, startMinute, "Once") // both params differ from the fixture inserted during setup
	if err != nil {
		t.Fatalf("encountered an error when inserting a new session: %v", err)
	}
	if newSessionId == fixtureWeeklySunday {
		t.Errorf("expected the new session ID to be different from the fixture (new: %v, fixture: %v)", newSessionId, fixtureWeeklySunday)
	}
	result, err := queries.GetSessionById(ctx, newSessionId)
	if err != nil {
		t.Fatalf("encountered an error when retrieving inserted record %v", newSessionId)
	}
	if !strings.HasPrefix(result.AddressFirstLine, string(houseNumber)) {
		t.Errorf("expected the address_first_line field to start with %v, got %v", houseNumber, result.AddressFirstLine)
	}

	if !(time.Now().UTC().Add(time.Minute*-1).Before(result.DtUpdatedUtc.Time) && time.Now().UTC().After(result.DtUpdatedUtc.Time)) {
		t.Errorf("expected dt_updated_utc to point to a value within the last minute before execution time of this statement, got %v (now: %v)", result.DtUpdatedUtc.Time, time.Now().UTC())
	}

	if result.Rating != 3 {
		t.Errorf("expected the rating to be the average of {1,5} = 3, got %v", result.Rating)
	}
}

func TestGetAllSessionsAsGeoJSON(t *testing.T) {
	var geojson types.FeatureCollection[types.SessionFeature[types.SessionProperties]]
	result, err := queries.GetAllSessionsAsGeoJSON(ctx)
	if err != nil {
		t.Fatalf("failed to retrieve sessions as geojson: %v", err)
	}
	err = json.Unmarshal(result, &geojson)
	if err != nil {
		t.Fatalf("failed to unmarshal json query result: %v", err)
	}
	if len(geojson.Features) == 0 {
		t.Error("expected at least one feature to be included in return")
	}
}

func TestGetSessionById(t *testing.T) {
	result, err := queries.GetSessionById(ctx, fixtureWeeklySunday)
	if err != nil {
		t.Fatalf("failed to retrieve session by id: %v", err)
	}
	if result.SessionID != fixtureWeeklySunday {
		t.Errorf("expected the returned ID (%v) to be the same as the fixture ID (%v)", result.SessionID, fixtureWeeklySunday)
	}
}

func TestGetSessionsByIdAsGeoJSON(t *testing.T) {
	var geojson types.SessionFeature[types.SessionProperties]
	result, err := queries.GetSessionByIdAsGeoJSON(ctx, fixtureWeeklySunday)
	if err != nil {
		t.Fatalf("failed to retrieve sessions as geojson: %v", err)
	}
	err = json.Unmarshal([]byte(result.(string)), &geojson)
	if err != nil {
		t.Fatalf("failed to unmarshal json query result: %v", err)
	}
}

func TestInsertAndRetrieveComment(t *testing.T) {
	insertedCommentId, err := queries.InsertSessionComment(ctx, InsertSessionCommentParams{
		Session: fixtureWeeklySunday,
		Author:  "foo",
		Content: "Example comment foo bar.",
	})
	if err != nil {
		t.Errorf("failed to insert session: %v", err)
	}

	result, err := queries.GetCommentsBySessionId(ctx, insertedCommentId)
	if err != nil {
		t.Errorf("failed to retrieve comments by session id: %v", err)
	}
	if result[0].Author != "foo" {
		t.Errorf("expected returned author to match the inserted comment record")
	}
}

func TestGetRatingsBySessionId(t *testing.T) {
	result, err := queries.GetRatingsBySessionId(ctx, fixtureWeeklySunday)
	if err != nil {
		t.Errorf("could not retrieve ratings by session ids: %v", err)
		t.FailNow()
	}
	if len(result) != 2 {
		t.Error("expected two ratings to be returned, instead got", len(result))
		t.FailNow()
	}

	if *result[0].Rating != 5 {
		t.Error("expected the first rating to be 5, instead got", *result[0].Rating)
	}

	if *result[1].Rating != 1 {
		t.Error("expected the second rating to be 1, instead got", *result[1].Rating)
	}
}

func TestGetSessionIdsByDate(t *testing.T) {
	result, err := queries.GetSessionIdsByDate(ctx, pgtype.Date{Time: time.Date(2024, 11, 19, 0, 0, 0, 0, time.UTC), Valid: true})
	if err != nil {
		t.Errorf("could not retrieve session ids by date: %v", err)
		t.FailNow()
	}
	var ids []int32
	var dates [][]any
	for idx := range result {
		s := result[idx].([]any)
		id := s[0].(int32)
		d := s[1].([]any)
		ids = append(ids, id)
		dates = append(dates, d)
	}
	if len(result) != 1 {
		t.Errorf("expected exactly 1 item in the result set, got ids %v", result)
		t.FailNow()
	}
	if len(dates[0]) != 1 {
		t.Errorf("expected the dates array to be of size 1, got %v", dates[0])
	}
	date := types.Date(dates[0][0].(time.Time))
	if ids[0] != fixtureDaily { // type cast
		t.Errorf("expected fixture 2 (%v), got %v", fixtureDaily, ids[0])
	}
	if date.String() != "2024-11-19" {
		t.Errorf("expected the dates attribute to be 2024-11-19, got %v", date.String())
	}
}

func TestGetSessionIdsByDate2(t *testing.T) {
	result, err := queries.GetSessionIdsByDate(ctx, pgtype.Date{Time: time.Date(2024, 11, 18, 0, 0, 0, 0, time.UTC), Valid: true})
	if err != nil {
		t.Errorf("could not retrieve session ids by date: %v", err)
		t.FailNow()
	}
	var ids []int32
	var dates [][]any
	for idx := range result {
		s := result[idx].([]any)
		id := s[0].(int32)
		d := s[1].([]any)
		ids = append(ids, id)
		dates = append(dates, d)
	}
	if len(result) != 2 {
		t.Errorf("expected exactly 2 item in the result set, got IDs %v", ids)
		t.FailNow()
	}
	if len(dates[1]) != 1 {
		t.Errorf("expected the dates array to be of size 1, got %v", dates[1])
	}
	date := types.Date(dates[0][0].(time.Time))
	if ids[1] != fixtureThirdOfMonthMonday { // type cast
		t.Errorf("expected fixture 3 (%v), got %v", fixtureThirdOfMonthMonday, ids[1])
	}
	if date.String() != "2024-11-18" {
		t.Errorf("expected the dates attribute to be 2024-11-18, got %v", date.String())
	}
}

func TestGetSessionIdsByDate3(t *testing.T) { // sunday
	result, err := queries.GetSessionIdsByDate(ctx, pgtype.Date{Time: time.Date(2024, 11, 30, 0, 0, 0, 0, time.UTC), Valid: true})
	if err != nil {
		t.Errorf("could not retrieve session ids by date: %v", err)
		t.FailNow()
	}
	if len(result) != 2 { // we expect id 2 (daily) and 4 (last saturday of the month)
		t.Error("expected exactly 2 item in the result set, got", result)
		t.FailNow()
	}
	for _, i := range result {
		s := i.([]any)
		sessionId := s[0].(int32)
		dates := s[1].([]any)
		if sessionId == fixtureDaily {
			if len(dates) != 1 { // daily
				t.Error("expected exactly 1 item in the date array, instead got", s[1])
			}
		} else if sessionId == fixtureLastOfMonthSaturday {
			// third monday of month, just once
			if len(dates) != 1 {
				t.Errorf("expected exactly 1 item in the date array for session %v, instead got %v", fixtureLastOfMonthSaturday, dates)
			}
			if !time.Date(2024, 11, 30, 0, 0, 0, 0, time.UTC).Equal(dates[0].(time.Time)) {
				t.Error("expected the date to be 2024-11-30, instead got", dates[0])
			}
		} else {
			t.Error("unexpected item", s)
		}
	}
}

func TestGetSessionIdsByDate4(t *testing.T) {
	result, err := queries.GetSessionIdsByDate(ctx, pgtype.Date{Time: time.Date(2024, 9, 11, 0, 0, 0, 0, time.UTC), Valid: true})
	if err != nil {
		t.Errorf("could not retrieve session ids by date: %v", err)
		t.FailNow()
	}
	var ids []int32
	var dates [][]any
	for idx := range result {
		s := result[idx].([]any)
		id := s[0].(int32)
		d := s[1].([]any)
		ids = append(ids, id)
		dates = append(dates, d)
	}
	if len(result) != 2 {
		t.Errorf("expected exactly 2 item in the result set, got IDs %v", ids)
		t.FailNow()
	}
	if len(dates[1]) != 1 {
		t.Errorf("expected the dates array to be of size 1, got %v", dates[1])
	}
	date := types.Date(dates[0][0].(time.Time))
	if ids[1] != fixtureFortnightlyWednesday { // type cast
		t.Errorf("expected fixture 5 (%v), got %v", fixtureFortnightlyWednesday, ids[1])
	}
	if date.String() != "2024-09-11" {
		t.Errorf("expected the dates attribute to be 2024-09-11, got %v", date.String())
	}
}

func TestGetSessionIdsByDate5(t *testing.T) {
	result, err := queries.GetSessionIdsByDate(ctx, pgtype.Date{Time: time.Date(2024, 11, 27, 0, 0, 0, 0, time.UTC), Valid: true})
	if err != nil {
		t.Errorf("could not retrieve session ids by date: %v", err)
		t.FailNow()
	}
	var ids []int32
	var dates [][]any
	for idx := range result {
		s := result[idx].([]any)
		id := s[0].(int32)
		d := s[1].([]any)
		ids = append(ids, id)
		dates = append(dates, d)
	}
	if len(result) != 1 {
		t.Errorf("expected exactly 1 item in the result set, got IDs %v", ids)
		t.FailNow()
	}
	if len(dates[0]) != 1 {
		t.Errorf("expected the dates array to be of size 1, got %v", dates[0])
	}
	date := types.Date(dates[0][0].(time.Time)) // type cast
	if ids[0] != fixtureDaily {
		t.Errorf("expected fixture 2 (%v), got %v", fixtureDaily, ids[0])
	}
	if date.String() != "2024-11-27" {
		t.Errorf("expected the dates attribute to be 2024-11-27, got %v", date.String())
	}
}

func TestGetSessionsByDate(t *testing.T) {
	result, err := queries.GetSessionsByDateAsGeoJSON(ctx, pgtype.Date{Time: time.Date(2024, 8, 20, 0, 0, 0, 0, time.UTC), Valid: true})
	if err != nil {
		t.Errorf("could not retrieve session ids by date: %v", err)
		t.FailNow()
	}
	var j types.SessionFeatureCollection
	if err := json.Unmarshal(result, &j); err != nil {
		t.Error("could not unmarshal:", err)
	}
	var ids []int32 = make([]int32, len(j.Features))
	for idx := range j.Features {
		ids[idx] = *j.Features[idx].Properties.SessionID
	}
	if len(j.Features) != 1 {
		t.Errorf("expected exactly 1 item in the result set, got ids %v", ids)
		t.FailNow()
	}
	if *j.Features[0].Properties.SessionID != fixtureDaily {
		t.Errorf("expected fixture 2 (%v), got %v", fixtureDaily, *j.Features[0].Properties.SessionID)
	}
	if j.Features[0].Properties.Dates == nil {
		t.Error("dates property shouldn't be nil")
	}
}

func TestGetSessionIdsByDateRange(t *testing.T) {
	result, err := queries.GetSessionIdsByDateRange(ctx, GetSessionIdsByDateRangeParams{
		StartDate: pgtype.Date{Time: time.Date(2024, 11, 16, 0, 0, 0, 0, time.UTC), Valid: true},
		EndDate:   pgtype.Date{Time: time.Date(2024, 11, 22, 0, 0, 0, 0, time.UTC), Valid: true},
	})
	if err != nil {
		t.Errorf("could not retrieve session ids by date range: %v", err)
		t.FailNow()
	}
	if len(result) != 4 {
		t.Error("expected exactly 3 items in the result set, instead got", result)
		t.FailNow()
	}
	for _, i := range result {
		s := i.([]any)
		sessionId := s[0].(int32)
		dates := s[1].([]any)
		if sessionId == fixtureDaily {
			if len(dates) != 7 { // daily, so must be the number of days between start and end
				t.Error("expected exactly 7 items in the date array, instead got", s[1])
			}
		} else if sessionId == fixtureThirdOfMonthMonday {
			// third monday of month, just once
			if len(dates) != 1 {
				t.Errorf("expected exactly 1 item in the date array for session %v, instead got %v", fixtureThirdOfMonthMonday, dates)
			}
			if !time.Date(2024, 11, 18, 0, 0, 0, 0, time.UTC).Equal(dates[0].(time.Time)) {
				t.Error("expected the date to be 2024-11-18, instead got", dates[0])
			}
		} else if sessionId == fixtureWeeklySunday {
			// every Sunday, once in this time window
			if len(dates) != 1 {
				t.Errorf("expected exactly 1 item in the date array for session %v, instead got %v", fixtureWeeklySunday, dates)

			}
			if !time.Date(2024, 11, 17, 0, 0, 0, 0, time.UTC).Equal(dates[0].(time.Time)) {
				t.Error("expected the date to be 2024-11-17, instead got", dates[0])
			}
		} else if sessionId == fixtureFortnightlyWednesday {
			if len(dates) != 1 {
				t.Errorf("expected exactly 1 item in the date array for session %v, instead got %v", fixtureFortnightlyWednesday, dates)
			}
			if !time.Date(2024, 11, 20, 0, 0, 0, 0, time.UTC).Equal(dates[0].(time.Time)) {
				t.Error("expected the date to be 2024-11-20, instead got", dates[0])
			}
		} else {
			t.Error("unexpected item", sessionId)
		}
	}
}

func TestGetSessionIdsByDateRange2(t *testing.T) {
	result, err := queries.GetSessionIdsByDateRange(ctx, GetSessionIdsByDateRangeParams{
		StartDate: pgtype.Date{Time: time.Date(2024, 11, 29, 0, 0, 0, 0, time.UTC), Valid: true},
		EndDate:   pgtype.Date{Time: time.Date(2024, 11, 30, 0, 0, 0, 0, time.UTC), Valid: true},
	})
	if err != nil {
		t.Errorf("could not retrieve session ids by date range: %v", err)
		t.FailNow()
	}
	if len(result) != 2 { // we expect id 2 (daily), 1 (weekly on a monday - 25th)
		t.Error("expected exactly 2 items in the result set, instead got", result)
		t.FailNow()
	}
	for _, i := range result {
		s := i.([]any)
		sessionId := s[0].(int32)
		dates := s[1].([]any)
		if sessionId == fixtureDaily {
			if len(dates) != 2 { // daily, so must be the number of days between start and end
				t.Error("expected exactly 2 items in the date array, instead got", s[1])
			}
			if !time.Date(2024, 11, 29, 0, 0, 0, 0, time.UTC).Equal(dates[0].(time.Time)) {
				t.Error("expected the date to be 2024-11-29, instead got", dates[0])
			}
		} else if sessionId == fixtureLastOfMonthSaturday {
			// third monday of month, just once
			if len(dates) != 1 {
				t.Errorf("expected exactly 1 item in the date array for session %v, instead got %v", fixtureLastOfMonthSaturday, dates)
			}
			if !time.Date(2024, 11, 30, 0, 0, 0, 0, time.UTC).Equal(dates[0].(time.Time)) {
				t.Error("expected the date to be 2024-11-30, instead got", dates[0])
			}
		} else {
			t.Error("unexpected item", sessionId)
		}
	}
}

func TestGetSessionsByDateRange(t *testing.T) {
	j, err := queries.GetSessionsByDateRangeAsGeoJSON(ctx, GetSessionsByDateRangeAsGeoJSONParams{
		StartDate: pgtype.Date{Time: time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC), Valid: true},
		EndDate:   pgtype.Date{Time: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), Valid: true},
	})
	if err != nil {
		t.Errorf("failed to retrieve session by date range: %v", err)
	}

	var result types.SessionWithVenueFeatureCollection
	if err := json.Unmarshal(j, &result); err != nil {
		t.Errorf("error when unmarshalling: %v", err)
		t.FailNow()
	}
	if len(result.Features) == 0 {
		t.Errorf("expected at least 1 feature, got %v", len(result.Features))
		t.FailNow()
	}
	if *result.Features[0].Properties.SessionID != fixtureWeeklySunday {
		t.Errorf("expected returned session (%v) to match the inserted fixture (%v)", *result.Features[0].Properties.SessionID, fixtureWeeklySunday)
	}
	if result.Features[0].Properties.Dates == nil {
		t.Error("dates property shouldn't be nil")
	}
}
