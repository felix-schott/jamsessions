package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	dbutils "github.com/felix-schott/jamsessions/backend/internal/db"
	migrationutils "github.com/felix-schott/jamsessions/backend/internal/migrations"
	"github.com/felix-schott/jamsessions/backend/internal/types"
	"github.com/jackc/pgx/v5/pgtype"
	geom "github.com/twpayne/go-geom"
)

// https://github.com/golang/go/issues/63309
func ptr[T any](t T) *T { return &t }

func TestCli(t *testing.T) {

	migrationsScript := "../../scripts/run-migrations.sh"

	// setup database connection
	pool, err := dbutils.CreatePool(ctx)
	if err != nil {
		log.Fatal("could not establish db connection:", err)
	}
	defer pool.Close()

	queries = dbutils.New(pool)
	var ctx = context.Background()

	// add test data - venue (tests use ephemeral databases so no need for cleanup after tests)
	testVenueId, err := queries.InsertVenue(ctx, dbutils.InsertVenueParams{
		VenueName:        "John Doe's Jazz Hole",
		AddressFirstLine: "11 Downing Street",
		City:             "London",
		Postcode:         "SW1A 2AA",
		Geom:             geom.NewPoint(geom.XY).MustSetCoords([]float64{-0.132, 51.513}),
		VenueWebsite:     ptr("https://www.test.com/"),
		Backline:         []string{"PA", "Drums"},
		VenueComments:    []string{"Comment 1", "comment 2"},
	})
	if err != nil {
		t.Errorf("the following error occurred when trying to insert a Venue record: %v", err)
		t.FailNow()
	}

	testSessionId, err := queries.InsertJamSession(ctx, dbutils.InsertJamSessionParams{
		SessionName:     "TEST_SESSION_12345",
		Venue:           testVenueId,
		Description:     "...",
		StartTimeUtc:    pgtype.Timestamptz{Time: time.Date(2024, 5, 6, 7, 6, 5, 4, time.UTC), Valid: true},
		DurationMinutes: 30,
		Interval:        "Weekly",
	})
	if err != nil {
		t.Errorf("the following error occurred when trying to insert a session record: %v", err)
		t.FailNow()
	}

	t.Run("UpdateVenueBackline", func(t *testing.T) {
		// temporary directory for testing
		migrationsDirectory := t.TempDir()
		migrationsArchive := filepath.Join(migrationsDirectory, "/archive")

		// write to migrations directory - cli call to add "Guitar_Amp" to backline field
		if fp, err := migrationutils.WriteMigration(fmt.Sprintf(`dbcli update venue %v "{"backline": ["PA", "Drums", "Guitar_Amp"]}"`, testVenueId), "test_update_venue", migrationsDirectory); err != nil {
			t.Errorf("could not write to file %v: %v", fp, err)
		}

		// simulate the manual execution of the script - note that if there were multiple tests, each test should have a separate migrationsDirectory for isolation
		var stderr bytes.Buffer
		cmd := exec.Command("bash", migrationsScript, "-y")
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "MIGRATIONS_DIRECTORY="+migrationsDirectory)
		cmd.Env = append(cmd.Env, "MIGRATIONS_ARCHIVE="+migrationsArchive)
		cmd.Stderr = &stderr
		cmd.Stdout = os.Stdout

		if err := cmd.Run(); err != nil {
			t.Errorf("an error occured when running migrations: %v: %v", err, stderr.String())
		}

		// check database for results
		record, err := queries.GetVenueById(ctx, testVenueId)
		if err != nil {
			t.Errorf("error when retrieving venue record (test fixture): %v", err)
		}
		if len(record.Backline) != 3 || record.Backline[0] != "PA" || record.Backline[1] != "Drums" || record.Backline[2] != "Guitar_Amp" {
			t.Errorf("expected the backline column to include the new addition. instead got: %v", record.Backline)
		}
	})

	t.Run("InsertVenueAndSession", func(t *testing.T) {
		// temporary directory for testing
		migrationsDirectory := t.TempDir()
		migrationsArchive := filepath.Join(migrationsDirectory, "/archive")

		// write to migrations directory - cli call to add "Guitar_Amp" to backline field
		testVenueProps := types.VenueProperties{
			VenueName:        ptr("Foo's Bar"),
			AddressFirstLine: ptr("10 Downing Street"),
			City:             ptr("London"),
			Postcode:         ptr("SW1A 2AA"),
			VenueWebsite:     ptr("foobar"),
		}
		testVenueJson, err := json.Marshal(testVenueProps)
		if err != nil {
			t.Errorf("could not marshal to json: %v", err)
			t.FailNow()
		}
		testSessionProps := types.SessionProperties{
			SessionName:     ptr("Foo's Session"),
			Venue:           ptr(int32(999999)), // replace this after serialisation
			Description:     ptr("Bla bla"),
			StartTimeUtc:    ptr(time.Date(2024, 5, 7, 1, 1, 1, 1, time.UTC)),
			DurationMinutes: ptr(int16(30)),
			Interval:        ptr("Weekly"),
		}
		testSessionJson, err := json.Marshal(testSessionProps)
		if err != nil {
			t.Errorf("could not marshal to json: %v", err)
			t.FailNow()
		}
		testSessionJson = []byte(strings.Replace(string(testSessionJson), "999999", "$new_id", -1)) // set venue to bash variable that will be evaluated  when the script runs
		if _, err := migrationutils.WriteMigration(fmt.Sprintf(`new_id=$(dbcli insert venue "%s");`+"\n"+`dbcli insert session "%s";`, testVenueJson, testSessionJson), "test_insert_session", migrationsDirectory); err != nil {
			t.Errorf("could not write migration: %v", err)
		}

		// simulate the manual execution of the script
		var stderr bytes.Buffer
		var stdout bytes.Buffer
		cmd := exec.Command("bash", migrationsScript, "-y")
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "MIGRATIONS_DIRECTORY="+migrationsDirectory)
		cmd.Env = append(cmd.Env, "MIGRATIONS_ARCHIVE="+migrationsArchive)
		cmd.Stderr = &stderr
		cmd.Stdout = &stdout

		if err := cmd.Run(); err != nil {
			if strings.Contains(stderr.String(), "status: 403") {
				log.Println("Nominatim unavailable (CI runs blocked, skip test)")
				return
			}
			t.Errorf("an error occured when running migrations: %v: %v", err, stderr.String())
		}

		// log.Println(stderr.String())

		// check database for results
		record, err := queries.GetVenueByName(ctx, "Foo's Bar")
		if err != nil {
			t.Errorf("error when retrieving inserted venue record: %v", err)
			t.FailNow()
		}
		if record.VenueName != "Foo's Bar" {
			t.Errorf("name (%v) doesn't match Foo's Bar", record.VenueName)
		}
		sessionId, err := strconv.Atoi(strings.ReplaceAll(stdout.String(), "\n", ""))
		log.Println("ID obtained from stdout:", sessionId)
		if err != nil {
			t.Errorf("could not parse stdout as session id: %v", err)
		}

		// check db
		rec, err := queries.GetSessionById(ctx, int32(sessionId))
		if err != nil {
			t.Error("error when retrieving inserted session record:", err)
		}
		if rec.SessionName != "Foo's Session" {
			t.Errorf("name of inserted session (%v) doesn't match Foo's Session", rec.SessionName)
		}
	})

	t.Run("InsertComment", func(t *testing.T) {
		// temporary directory for testing
		migrationsDirectory := t.TempDir()
		migrationsArchive := filepath.Join(migrationsDirectory, "/archive")

		testCommentJson, err := json.Marshal(dbutils.InsertSessionCommentParams{
			Session: testSessionId,
			Author:  "test author",
			Content: "This is a comment.",
		})
		if err != nil {
			t.Error("Couldn't marshal json", err)
			t.FailNow()
		}

		if fp, err := migrationutils.WriteMigration(fmt.Sprintf(`dbcli insert comment "%s";`, testCommentJson), "test_insert_comment", migrationsDirectory); err != nil {
			t.Errorf("could not write to file %v: %v", fp, err)
		}

		// run migrations
		var stderr bytes.Buffer
		var stdout bytes.Buffer
		cmd := exec.Command("bash", migrationsScript, "-y")
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "MIGRATIONS_DIRECTORY="+migrationsDirectory)
		cmd.Env = append(cmd.Env, "MIGRATIONS_ARCHIVE="+migrationsArchive)
		cmd.Stderr = &stderr
		cmd.Stdout = &stdout

		if err := cmd.Run(); err != nil {
			t.Errorf("an error occured when running migrations: %v: %v", err, stderr.String())
		}

		// check db
		rec, err := queries.GetCommentsBySessionId(ctx, testSessionId)
		if err != nil {
			t.Error("error when retrieving inserted session record:", err)
		}
		if rec[0].Author != "test author" {
			t.Errorf("name of inserted comment author (%v) doesn't match 'test author'", rec[0].Author)
		}
	})
}
