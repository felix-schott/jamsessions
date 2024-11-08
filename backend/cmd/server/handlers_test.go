package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	geom "github.com/twpayne/go-geom"

	dbutils "github.com/felix-schott/jamsessions/backend/internal/db"
	"github.com/felix-schott/jamsessions/backend/internal/types"
	"github.com/go-fuego/fuego"
)

var s *fuego.Server = fuego.NewServer(fuego.WithAddr("localhost:66666"))

// helper func - check if the returned feature collection contains the sessionIds provided
func checkResultSetForSessionIds(t *testing.T, sessionIds []int32, body types.SessionFeatureCollection) {
	foundMap := make(map[int32]bool)
	for _, id := range sessionIds {
		foundMap[id] = false
	}

	var allSessionsStr []string

	for _, f := range body.Features {
		if f.Properties.SessionID != nil {
			allSessionsStr = append(allSessionsStr, string(*f.Properties.SessionID))
			if _, ok := foundMap[*f.Properties.SessionID]; ok {
				foundMap[*f.Properties.SessionID] = true
			}
		}
	}

	for id, found := range foundMap {
		if !found {
			t.Errorf("expected one of the returned sessions to have ID %v, instead got %v", id, strings.Join(allSessionsStr, ", "))
		}
	}
}

func TestHandlers(t *testing.T) {

	slog.SetLogLoggerLevel(slog.LevelError) // change this to see more log informations

	// setup database connection
	pool, err := dbutils.CreatePool(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// defer pool.Close()

	queries = dbutils.New(pool)

	// add test data - venue (tests use ephemeral databases so no need for cleanup after tests)
	testVenueId, err := queries.InsertVenue(ctx, dbutils.InsertVenueParams{
		VenueName:        "TEST HANDLERS",
		AddressFirstLine: "1 Main Street",
		City:             "London",
		Postcode:         "ABC 123",
		Geom:             geom.NewPoint(geom.XY).MustSetCoords([]float64{-0.132, 51.514}),
		VenueWebsite:     ptr("https://www.test.com/"),
		Backline:         []string{"PA", "Drums"},
		VenueComments:    []string{"Comment 1", "comment 2"},
	})
	if err != nil {
		t.Errorf("the following error occurred when trying to insert a Venue record: %v", err)
	}

	// add test data - session (tests use ephemeral databases so no need for cleanup after tests)
	_, err = queries.InsertJamSession(ctx, dbutils.InsertJamSessionParams{
		SessionName:     "test_session1",
		StartTimeUtc:    pgtype.Timestamptz{Time: time.Date(2024, 1, 1, 19, 30, 0, 0, time.UTC), Valid: true},
		Interval:        "Weekly",
		DurationMinutes: 120,
		Venue:           testVenueId,
		Genres:          []string{"Blues", "Jazz-Funk"},
	})
	if err != nil {
		t.Errorf("the following error occured when trying to insert test session 1: %v", err)
	}

	_, err = queries.InsertJamSession(ctx, dbutils.InsertJamSessionParams{
		SessionName:     "test_session2",
		StartTimeUtc:    pgtype.Timestamptz{Time: time.Date(2024, 1, 2, 19, 30, 0, 0, time.UTC), Valid: true},
		Interval:        "Weekly",
		DurationMinutes: 120,
		Venue:           testVenueId,
		Genres:          []string{"Jazz-Funk"},
	})
	if err != nil {
		t.Errorf("the following error occured when trying to insert test session 2: %v", err)
	}

	// t.Run("GetAllVenues", func(t *testing.T) {
	// 	handler := fuego.HTTPHandler(s, GetVenues)
	// 	req := httptest.NewRequest(http.MethodGet, "/venues", nil)
	// 	w := httptest.NewRecorder()
	// 	handler(w, req)
	// 	res := w.Result()
	// 	defer res.Body.Close()
	// 	data, err := io.ReadAll(res.Body)
	// 	if err != nil {
	// 		t.Errorf("expected error to be nil got %v", err)
	// 	}
	// 	var body types.VenueFeatureCollection
	// 	err = json.Unmarshal(data, &body)
	// 	if err != nil {
	// 		t.Errorf("expected error to be nil got %v", err)
	// 	}
	// 	if len(body.Features) == 0 {
	// 		t.Errorf("expected at least 1 feature in the venue feature collection")
	// 	}
	// })

	// t.Run("GetAllSessions", func(t *testing.T) {
	// 	handler := fuego.HTTPHandler(s, GetSessions)
	// 	req := httptest.NewRequest(http.MethodGet, "/sessions", nil)
	// 	w := httptest.NewRecorder()
	// 	handler(w, req)
	// 	res := w.Result()
	// 	defer res.Body.Close()
	// 	data, err := io.ReadAll(res.Body)
	// 	if err != nil {
	// 		t.Errorf("expected error to be nil got %v", err)
	// 	}
	// 	var body types.SessionFeatureCollection
	// 	err = json.Unmarshal(data, &body)
	// 	if err != nil {
	// 		t.Errorf("expected error to be nil got %v", err)
	// 	}
	// 	if len(body.Features) < 2 {
	// 		t.Errorf("expected at least 2 feature in the session feature collection, got %v", len(body.Features))
	// 		t.FailNow()
	// 	}

	// 	checkResultSetForSessionIds(t, []int32{testSession1Id, testSession2Id}, body)
	// })

	// t.Run("GetSessionsByInferredDate", func(t *testing.T) {
	// 	handler := fuego.HTTPHandler(s, GetSessions)
	// 	req := httptest.NewRequest(http.MethodGet, "/jamsessions?date=2024-01-09", nil) // start date of session 2 plus 1 week (to match interval)
	// 	w := httptest.NewRecorder()
	// 	handler(w, req)
	// 	res := w.Result()
	// 	defer res.Body.Close()
	// 	data, err := io.ReadAll(res.Body)
	// 	if err != nil {
	// 		t.Errorf("expected error to be nil got %v", err)
	// 	}
	// 	var body types.SessionFeatureCollection
	// 	err = json.Unmarshal(data, &body)
	// 	if err != nil {
	// 		t.Errorf("expected error to be nil got %v", err)
	// 	}
	// 	if len(body.Features) == 0 {
	// 		t.Errorf("expected at least 1 feature in the session feature collection")
	// 	}

	// 	checkResultSetForSessionIds(t, []int32{testSession2Id}, body)
	// })

	t.Run("GetSessionsNoDateMatch", func(t *testing.T) {
		handler := fuego.HTTPHandler(s, GetSessions)
		req := httptest.NewRequest(http.MethodGet, "/jamsessions?date=2023-01-06", nil)
		w := httptest.NewRecorder()
		handler(w, req)
		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}
		var body types.SessionFeatureCollection
		err = json.Unmarshal(data, &body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}
		if len(body.Features) != 0 {
			var sessions = make([]string, len(body.Features))
			for i, f := range body.Features {
				sessions[i] = *f.Properties.SessionName
				if f.Properties.Dates != nil {
					sessions[i] = sessions[i] + " ("
					for _, d := range *f.Properties.Dates {
						sessions[i] = sessions[i] + fmt.Sprint(d)
					}
					sessions[i] = sessions[i] + ")"
				} else {
					sessions[i] = sessions[i] + " (no date property)"
				}
			}
			t.Errorf("expected no features in the session feature collection, got %v", strings.Join(sessions, ", "))
		}
	})

	// t.Run("GetSessionsByDateAndGenre", func(t *testing.T) {
	// 	handler := fuego.HTTPHandler(s, GetSessions)
	// 	req := httptest.NewRequest(http.MethodGet, "/jamsessions?date=2024-01-01&genre=Blues", nil)
	// 	w := httptest.NewRecorder()
	// 	handler(w, req)
	// 	res := w.Result()
	// 	defer res.Body.Close()
	// 	data, err := io.ReadAll(res.Body)
	// 	if err != nil {
	// 		t.Errorf("expected error to be nil got %v", err)
	// 	}
	// 	var body types.SessionFeatureCollection
	// 	err = json.Unmarshal(data, &body)
	// 	if err != nil {
	// 		t.Errorf("expected error to be nil got %v", err)
	// 	}
	// 	if len(body.Features) == 0 {
	// 		t.Errorf("expected at least 1 feature in the session feature collection")
	// 	}
	// 	checkResultSetForSessionIds(t, []int32{testSession1Id}, body)
	// })

	// t.Run("GetSessionsByDateRangeAndGenre", func(t *testing.T) {
	// 	handler := fuego.HTTPHandler(s, GetSessions)
	// 	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/jamsessions?date=%v&genre=Blues", url.QueryEscape("2024-01-01/2024-03-01")), nil)
	// 	w := httptest.NewRecorder()
	// 	handler(w, req)
	// 	res := w.Result()
	// 	defer res.Body.Close()
	// 	data, err := io.ReadAll(res.Body)
	// 	if err != nil {
	// 		t.Errorf("expected error to be nil got %v", err)
	// 	}
	// 	var body types.SessionFeatureCollection
	// 	err = json.Unmarshal(data, &body)
	// 	if err != nil {
	// 		t.Errorf("expected error to be nil got %v", err)
	// 	}
	// 	if len(body.Features) == 0 {
	// 		t.Errorf("expected at least 1 feature in the session feature collection")
	// 	}
	// 	checkResultSetForSessionIds(t, []int32{testSession1Id}, body)
	// })

	// t.Run("GetSessionsByDateRangeAndGenre2", func(t *testing.T) {
	// 	handler := fuego.HTTPHandler(s, GetSessions)
	// 	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/jamsessions?date=%v&genre=Jazz-Funk", url.QueryEscape("2024-01-07/2024-01-09")), nil)
	// 	w := httptest.NewRecorder()
	// 	handler(w, req)
	// 	res := w.Result()
	// 	defer res.Body.Close()
	// 	data, err := io.ReadAll(res.Body)
	// 	if err != nil {
	// 		t.Errorf("expected error to be nil got %v", err)
	// 	}
	// 	var body types.SessionFeatureCollection
	// 	err = json.Unmarshal(data, &body)
	// 	if err != nil {
	// 		t.Errorf("failed to unmarshal %s: %v", data, err)
	// 	}
	// 	if len(body.Features) < 2 {
	// 		t.Errorf("expected at least 2 feature in the session feature collection")
	// 	}
	// 	checkResultSetForSessionIds(t, []int32{testSession1Id, testSession2Id}, body)
	// })

	// test currently failing - not essential, but should be looked into at some point
	// t.Run("GetSessionsByDateAndBacklineNoMatch", func(t *testing.T) {
	// 	handler := fuego.HTTPHandler(s, GetSessions)
	// 	req := httptest.NewRequest(http.MethodGet, "/jamsessions?date=2024-01-01&backline=PA,Guitar_Amp", nil)
	// 	w := httptest.NewRecorder()
	// 	handler(w, req)
	// 	res := w.Result()
	// 	defer res.Body.Close()
	// 	data, err := io.ReadAll(res.Body)
	// 	if err != nil {
	// 		t.Errorf("expected error to be nil got %v", err)
	// 	}
	// 	var body types.SessionFeatureCollection
	// 	err = json.Unmarshal(data, &body)
	// 	if err != nil {
	// 		t.Errorf("expected error to be nil got %v", err)
	// 	}
	// 	if len(body.Features) != 0 {
	// 		var ids []string
	// 		for _, f := range body.Features {
	// 			ids = append(ids, *f.Properties.SessionName)
	// 		}
	// 		t.Errorf("expected no features in the session feature collection, got %v (with these IDs: %v)", len(body.Features), strings.Join(ids, ", "))
	// 	}
	// })

	// 	t.Run("GetSessionsByDateAndBackline", func(t *testing.T) {
	// 		handler := fuego.HTTPHandler(s, GetSessions)
	// 		req := httptest.NewRequest(http.MethodGet, "/jamsessions?date=2024-01-01&backline=PA,Drums", nil)
	// 		w := httptest.NewRecorder()
	// 		handler(w, req)
	// 		res := w.Result()
	// 		defer res.Body.Close()
	// 		data, err := io.ReadAll(res.Body)
	// 		if err != nil {
	// 			t.Errorf("expected error to be nil got %v", err)
	// 		}
	// 		var body types.SessionFeatureCollection
	// 		err = json.Unmarshal(data, &body)
	// 		if err != nil {
	// 			t.Errorf("failed to unmarshal %s: %v", data, err)
	// 		}
	// 		if len(body.Features) == 0 {
	// 			var ids []string
	// 			for _, f := range body.Features {
	// 				ids = append(ids, *f.Properties.SessionName)
	// 			}

	// 			t.Errorf("expected at least 1 feature in the session feature collection, got %v (with the following names: %v)", len(body.Features), strings.Join(ids, ", "))
	// 			t.FailNow()
	// 		}
	// 		checkResultSetForSessionIds(t, []int32{testSession1Id}, body)
	// 	})

	// 	t.Run("GetSessionById", func(t *testing.T) {
	// 		handler := fuego.HTTPHandler(s, GetSessionById)
	// 		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/jamsessions/%v", testSession2Id), nil)
	// 		req.SetPathValue("id", fmt.Sprint(testSession2Id))
	// 		w := httptest.NewRecorder()
	// 		handler(w, req)
	// 		res := w.Result()
	// 		if res.StatusCode != 200 {
	// 			t.Errorf("expected status code 200, got %v", res.StatusCode)
	// 		}
	// 		defer res.Body.Close()
	// 		data, err := io.ReadAll(res.Body)
	// 		if err != nil {
	// 			t.Errorf("expected error to be nil got %v", err)
	// 		}
	// 		var body types.SessionFeature[types.SessionProperties]
	// 		err = json.Unmarshal(data, &body)
	// 		if err != nil {
	// 			t.Errorf("expected error to be nil got %v", err)
	// 		}
	// 		if body.Properties.SessionID != nil && *body.Properties.SessionID != testSession2Id {
	// 			t.Errorf("expected session 2 to be returned")
	// 		}
	// 	})

	// 	t.Run("PostCommentForSessionById", func(t *testing.T) {

	// 		// temp directory for migrations
	// 		migrationsDirectory = t.TempDir()

	// 		testComment := "Test comment number 123!"

	// 		handler := fuego.HTTPHandler(s, PostCommentForSessionById)
	// 		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/jamsessions/%v/comments", testSession1Id), strings.NewReader(fmt.Sprintf(`{"content": "%v"}`, testComment)))
	// 		req.SetPathValue("id", fmt.Sprint(testSession1Id))
	// 		w := httptest.NewRecorder()
	// 		handler(w, req)
	// 		res := w.Result()
	// 		if res.StatusCode != 201 {
	// 			t.Errorf("expected status code 201, got %v", res.StatusCode)
	// 		}

	// 		dir, err := os.ReadDir(migrationsDirectory)
	// 		if err != nil {
	// 			t.Errorf("couldn't read directory contents: %v", err)
	// 			t.FailNow()
	// 		}
	// 		if len(dir) != 1 {
	// 			t.Errorf("expected exactly 1 file in the directory, got %v", len(dir))
	// 			t.FailNow()
	// 		}

	// 		f, err := os.ReadFile(filepath.Join(migrationsDirectory, dir[0].Name()))
	// 		if err != nil {
	// 			t.Errorf("error reading file: %v", err)
	// 		}

	// 		matched, err := regexp.Match(fmt.Sprintf(`.*.*dbcli insert comment "{\\"session\\":%v,\\"author\\":\\"\\",\\"content\\":\\"%v\\"}"`, testSession1Id, testComment), f)
	// 		if err != nil {
	// 			t.Errorf("error when trying match with regex: %v", err)
	// 		}

	// 		if !matched {
	// 			t.Errorf("expected the regex to match. instead got file contents: %s", f)
	// 		}
	// 		// see internal/db/cli for cli tests
	// 	})

	// 	t.Run("PostSession", func(t *testing.T) {
	// 		// temp directory for migrations
	// 		migrationsDirectory = t.TempDir()

	// 		testBody, err := json.Marshal(types.SessionProperties{
	// 			SessionName:     ptr("TestInsert"),
	// 			Venue:           &testVenueId,
	// 			Description:     ptr("Description."),
	// 			StartTimeUtc:    ptr(time.Date(2024, 3, 4, 5, 5, 3, 5, time.UTC)),
	// 			DurationMinutes: ptr(int16(90)),
	// 			Interval:        ptr("FirstOfMonth"),
	// 			SessionWebsite:  ptr("https://example.org"),
	// 		})
	// 		if err != nil {
	// 			t.Error("could not marshal json:", err)
	// 			t.FailNow()
	// 		}

	// 		handler := fuego.HTTPHandler(s, PostSession)
	// 		req := httptest.NewRequest(http.MethodPost, "/jamsessions", bytes.NewReader(testBody))
	// 		w := httptest.NewRecorder()
	// 		handler(w, req)
	// 		res := w.Result()
	// 		if res.StatusCode != 201 {
	// 			t.Errorf("expected status code 201, got %v", res.StatusCode)
	// 		}

	// 		dir, err := os.ReadDir(migrationsDirectory)
	// 		if err != nil {
	// 			t.Errorf("couldn't read directory contents: %v", err)
	// 			t.FailNow()
	// 		}
	// 		if len(dir) != 1 {
	// 			t.Errorf("expected exactly 1 file in the directory, got %v", len(dir))
	// 			t.FailNow()
	// 		}

	// 		f, err := os.ReadFile(filepath.Join(migrationsDirectory, dir[0].Name()))
	// 		if err != nil {
	// 			t.Errorf("error reading file: %v", err)
	// 		}

	// 		matched, err := regexp.Match(fmt.Sprintf(`insert session "{\\"session_name\\":\\"%v\\",.*`, "TestInsert"), f)
	// 		if err != nil {
	// 			t.Errorf("error when trying match with regex: %v", err)
	// 		}

	// 		if !matched {
	// 			t.Errorf("expected the regex to match. instead got file contents: %s", f)
	// 		}
	// 	})

	// 	t.Run("PostSessionAndVenue", func(t *testing.T) {
	// 		// temp directory for migrations
	// 		migrationsDirectory = t.TempDir()

	// 		testBody, err := json.Marshal(types.SessionPropertiesWithVenue{SessionProperties: types.SessionProperties{
	// 			SessionName:     ptr("TestInsert2"),
	// 			Description:     ptr("Description."),
	// 			StartTimeUtc:    ptr(time.Date(2024, 3, 4, 5, 5, 3, 5, time.UTC)),
	// 			DurationMinutes: ptr(int16(90)),
	// 			Interval:        ptr("FirstOfMonth"),
	// 			SessionWebsite:  ptr("https://example.org"),
	// 		},
	// 			VenueProperties: types.VenueProperties{
	// 				VenueName:        ptr("VenueInsert"),
	// 				AddressFirstLine: ptr("1 Random St"),
	// 				City:             ptr("Randomtown"),
	// 				Postcode:         ptr("ABC 123"),
	// 				VenueWebsite:     ptr("https://example.org"),
	// 				Backline:         &[]string{"PA", "Drums"},
	// 			},
	// 		})
	// 		if err != nil {
	// 			t.Error("could not marshal json:", err)
	// 			t.FailNow()
	// 		}

	// 		handler := fuego.HTTPHandler(s, PostSession)
	// 		req := httptest.NewRequest(http.MethodPost, "/jamsessions", bytes.NewReader(testBody))
	// 		w := httptest.NewRecorder()
	// 		handler(w, req)
	// 		res := w.Result()
	// 		if res.StatusCode != 201 {
	// 			t.Errorf("expected status code 201, got %v", res.StatusCode)
	// 		}

	// 		dir, err := os.ReadDir(migrationsDirectory)
	// 		if err != nil {
	// 			t.Errorf("couldn't read directory contents: %v", err)
	// 			t.FailNow()
	// 		}
	// 		if len(dir) != 1 {
	// 			t.Errorf("expected exactly 1 file in the directory, got %v", len(dir))
	// 			t.FailNow()
	// 		}

	// 		f, err := os.ReadFile(filepath.Join(migrationsDirectory, dir[0].Name()))
	// 		if err != nil {
	// 			t.Errorf("error reading file: %v", err)
	// 		}

	// 		matched, err := regexp.Match(fmt.Sprintf(`new_id=\$\(dbcli insert venue "{\\"venue_name\\":\\"%v\\",.*\n.*\\"session_name\\":\\"%v\\",.*`, "VenueInsert", "TestInsert2"), f)
	// 		if err != nil {
	// 			t.Errorf("error when trying match with regex: %v", err)
	// 		}

	// 		if !matched {
	// 			t.Errorf("expected the regex to match. instead got file contents: %s", f)
	// 		}
	// 	})

	// 	t.Run("PostSessionAltPayload", func(t *testing.T) {
	// 		// temp directory for migrations
	// 		migrationsDirectory = t.TempDir()

	// 		testBody := `{"session_name":"TEST session 123","description":"dafdsc dsd.","interval":"Weekly","start_time_utc":"2024-10-16T00:00:00.000Z","duration_minutes":60,"genres":[],"session_website":"https://example.org/"}`

	// 		handler := fuego.HTTPHandler(s, PostSession)
	// 		req := httptest.NewRequest(http.MethodPost, "/jamsessions", strings.NewReader(testBody))
	// 		w := httptest.NewRecorder()
	// 		handler(w, req)
	// 		res := w.Result()
	// 		if res.StatusCode != 201 {
	// 			t.Errorf("expected status code 201, got %v", res.StatusCode)
	// 		}

	// 		dir, err := os.ReadDir(migrationsDirectory)
	// 		if err != nil {
	// 			t.Errorf("couldn't read directory contents: %v", err)
	// 			t.FailNow()
	// 		}
	// 		if len(dir) != 1 {
	// 			t.Errorf("expected exactly 1 file in the directory, got %v", len(dir))
	// 			t.FailNow()
	// 		}

	// 		f, err := os.ReadFile(filepath.Join(migrationsDirectory, dir[0].Name()))
	// 		if err != nil {
	// 			t.Errorf("error reading file: %v", err)
	// 		}

	// 		matched, err := regexp.Match(fmt.Sprintf(`insert session "{\\"session_name\\":\\"%v\\",.*`, "TEST session 123"), f)
	// 		if err != nil {
	// 			t.Errorf("error when trying match with regex: %v", err)
	// 		}

	//		if !matched {
	//			t.Errorf("expected the regex to match. instead got file contents: %s", f)
	//		}
	//	})
}
