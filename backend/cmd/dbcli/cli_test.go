package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	dbutils "github.com/felix-schott/london-jam-sessions/internal/db"
	"github.com/felix-schott/london-jam-sessions/internal/types"
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
		VenueName:        "__TEST__",
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

	t.Run("UpdateVenueBackline", func(t *testing.T) {
		// temporary directory for testing
		migrationsDirectory := t.TempDir()
		migrationsArchive := filepath.Join(migrationsDirectory, "/archive")

		// write to migrations directory - cli call to add "Guitar_Amp" to backline field
		fp := filepath.Join(migrationsDirectory, "UpdateVenueBackline.sh")
		if err := os.WriteFile(fp, []byte(fmt.Sprintf(`dbcli update venue %v '{"backline": ["PA", "Drums", "Guitar_Amp"]}'`, testVenueId)), fs.FileMode(int(0755))); err != nil {
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
			VenueName:        ptr("TEST_VENUE"),
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
			SessionName:     ptr("TEST_SESSION"),
			Venue:           ptr(int32(999999)), // replace this after serialisation
			Description:     ptr("Bla bla"),
			StartTimeUtc:    ptr(time.Date(2024, 1, 1, 1, 1, 1, 1, time.UTC)),
			DurationMinutes: ptr(int16(30)),
			Interval:        ptr("Weekly"),
		}
		testSessionJson, err := json.Marshal(testSessionProps)
		if err != nil {
			t.Errorf("could not marshal to json: %v", err)
			t.FailNow()
		}
		testSessionJson = []byte(strings.Replace(string(testSessionJson), "999999", "$new_id", -1)) // set venue to bash variable that will be evaluated  when the script runs
		fp := filepath.Join(migrationsDirectory, "InsertVenueAndSession.sh")
		if err := os.WriteFile(fp, []byte(fmt.Sprintf(`new_id=$(dbcli insert venue '%v');`+"\n"+`dbcli insert session "%v";`, string(testVenueJson), strings.Replace(string(testSessionJson), `"`, `\"`, -1))), fs.FileMode(int(0755))); err != nil {
			t.Errorf("could not write to file %v: %v", fp, err)
		}

		// simulate the manual execution of the script - note that if there were multiple tests, each test should have a separate migrationsDirectory for isolation
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
		record, err := queries.GetVenueByName(ctx, "TEST_VENUE")
		if err != nil {
			t.Errorf("error when retrieving inserted venue record: %v", err)
			t.FailNow()
		}
		if record.VenueName != "TEST_VENUE" {
			t.Errorf("name (%v) doesn't match TEST_VENUE", record.VenueName)
		}
		sessionId, err := strconv.Atoi(strings.ReplaceAll(stdout.String(), "\n", ""))
		log.Println("ID obtained from stdout:", sessionId)
		if err != nil {
			t.Errorf("could not parse stdout as session id: %v", err)
		}
		rec, err := queries.GetSessionById(ctx, int32(sessionId))
		if err != nil {
			t.Errorf("error when retrieving inserted session record")
		}
		if rec.SessionName != "TEST_SESSION" {
			t.Errorf("name of inserted session (%v) doesn't match TEST_SESSION", rec.SessionName)
		}
	})
}
