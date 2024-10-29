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
var fixtureSessionId int32

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
		log.Fatalf("test expects a pre-populated DB. Tables 'jamsessions' and 'venues' found missing. Number of table: %v", numTables)
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

	// add one dummy record to both Location and JamSession table
	fixtureSessionId, err = insertSession("Ronnie Scott's Jazz Club", 12, 15)
	if err != nil {
		log.Fatal(err)
	}
}

func teardown() {
	if fixtureSessionId != 0 {
		err := queries.DeleteVenueByJamSessionId(ctx, fixtureSessionId)
		if err != nil {
			log.Fatalf("could not delete venue previously inserted for session id %v: %v", fixtureSessionId, err)
		}
		err = queries.DeleteJamSessionById(ctx, fixtureSessionId)
		if err != nil {
			log.Fatalf("could not delete session with id %v", fixtureSessionId)
		}
	}
}

// https://github.com/golang/go/issues/63309
func ptr[T any](t T) *T { return &t }

// Helper function that inserts a JamSession and corresponding Location record.
// The two params houseNumber and startMinute can be used to control uniqueness within the
// Location and JamSession table, respectively.
func insertSession(locationName string, houseNumber uint8, startMinute int) (int32, error) {
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
		StartTimeUtc:    pgtype.Timestamptz{Time: time.Date(2024, 1, 1, 19, startMinute, 0, 0, time.UTC), Valid: true},
		Interval:        "Weekly",
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
	newSessionId, err := insertSession("Ronnie Scott's Blues Club", houseNumber, startMinute) // both params differ from the fixture inserted during setup
	if err != nil {
		t.Fatalf("encountered an error when inserting a new session: %v", err)
	}
	if newSessionId == fixtureSessionId {
		t.Errorf("expected the new session ID to be different from the fixture (new: %v, fixture: %v)", newSessionId, fixtureSessionId)
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
	result, err := queries.GetSessionById(ctx, fixtureSessionId)
	if err != nil {
		t.Fatalf("failed to retrieve session by id: %v", err)
	}
	if result.SessionID != fixtureSessionId {
		t.Errorf("expected the returned ID (%v) to be the same as the fixture ID (%v)", result.SessionID, fixtureSessionId)
	}
}

func TestGetSessionsByIdAsGeoJSON(t *testing.T) {
	var geojson types.SessionFeature[types.SessionProperties]
	result, err := queries.GetSessionByIdAsGeoJSON(ctx, fixtureSessionId)
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
		Session: fixtureSessionId,
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
