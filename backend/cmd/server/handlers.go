package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	dbutils "github.com/felix-schott/jamsessions/backend/internal/db"
	types "github.com/felix-schott/jamsessions/backend/internal/types"
	"github.com/go-fuego/fuego"
	"github.com/jackc/pgx/v5/pgtype"
)

type Interval int

const (
	Once Interval = iota
	Daily
	Weekly
	Fortnightly
	FirstOfMonth
	SecondOfMonth
	ThirdOfMonth
	FourthOfMonth
	LastOfMonth
)

type Backline int

const (
	PA Backline = iota
	GuitarAmp
	BassAmp
	Drums
	Keys
	Mic
	MiscPercussion
)

var backlineStrToEnum = map[string]Backline{
	"PA":             PA,
	"Guitar_Amp":     GuitarAmp,
	"Bass_Amp":       BassAmp,
	"Drums":          Drums,
	"Microphone":     Mic,
	"MiscPercussion": MiscPercussion,
	"Keys":           Keys,
}

func (b Backline) String() string {
	for k, v := range backlineStrToEnum {
		if v == b {
			return k
		}
	}
	return ""
}

type Genre int

// TODO enum for genres
const (
	Any Genre = iota
	StraightAhead
	JazzFunk
	Fusion
	LatinJazz
	ModernJazz
	TradJazz
	Funk
	Blues
	Folk
	Rock
	WorldMusic
)

var genreStrToEnum = map[string]Genre{
	"Straight-Ahead_Jazz": StraightAhead,
	"Jazz-Funk":           JazzFunk,
	"Fusion":              Fusion,
	"Latin_Jazz":          LatinJazz,
	"Funk":                Funk,
	"Blues":               Blues,
	"Folk":                Folk,
	"Rock":                Rock,
	"World_Music":         WorldMusic,
	"Modern_Jazz":         ModernJazz,
	"Trad_Jazz":           TradJazz,
}

func (b Genre) String() string {
	for k, v := range genreStrToEnum {
		if v == b {
			return k
		}
	}
	return ""
}

// func matchesIrregularInterval(d *time.Time, i Interval) bool {
// 	nthInMonth := math.Ceil(float64(d.Day()) / 7)
// 	if i == FirstOfMonth && nthInMonth == 1 {
// 		return true
// 	}
// 	if i == SecondOfMonth && nthInMonth == 2 {
// 		return true
// 	}
// 	if i == ThirdOfMonth && nthInMonth == 3 {
// 		return true
// 	}
// 	if i == FourthOfMonth && nthInMonth == 4 {
// 		return true
// 	}
// 	if i == LastOfMonth && d.AddDate(0, 0, 7).Month() != d.Month() {
// 		return true
// 	}
// 	return false
// }

func GetVenues(c *fuego.ContextNoBody) (types.VenueFeatureCollection, error) {
	var geojson types.FeatureCollection[types.VenueFeature]
	result, err := queries.GetAllVenuesAsGeoJSON(ctx)
	if err != nil {
		return geojson, err
	}
	slog.Info("foo", "result", string(result))
	json.Unmarshal(result, &geojson)
	slog.Info("GetVenues", "result", geojson)
	return geojson, nil
}

func GetSessions(c *fuego.ContextNoBody) (types.SessionWithVenueFeatureCollection, error) {
	queryParams := c.QueryParams()
	var date time.Time
	var backline []string
	var genre string
	var geojson types.SessionWithVenueFeatureCollection

	invalidKeys := make([]string, 0, len(queryParams))
	i := 0
	for k := range queryParams {
		if k == "date" {
			dateParsed, err := time.Parse(time.DateOnly, c.QueryParam("date"))
			if err != nil {
				return geojson, fuego.BadRequestError{Detail: fmt.Sprintf("failed to parse %v as a date, please provide dates as 'YYYY-MM-DD'", c.QueryParam("date"))}
			}
			date = dateParsed
		} else if k == "backline" {
			backline = strings.Split(c.QueryParam("backline"), ",")
			for _, b := range backline {
				_, ok := backlineStrToEnum[b]
				if !ok {
					return geojson, fuego.BadRequestError{Detail: fmt.Sprintf("%v is not a valid value for 'backline'", b)}
				}
			}
		} else if k == "genre" {
			genre = c.QueryParam("genre")
			_, ok := genreStrToEnum[genre]
			if !ok {
				return geojson, fuego.BadRequestError{Detail: fmt.Sprintf("%v is not a valid value for 'genre'", genre)}
			}
		} else {
			invalidKeys = append(invalidKeys, k)
		}
		i++
	}
	if len(invalidKeys) != 0 {
		return types.SessionWithVenueFeatureCollection{}, fuego.BadRequestError{Detail: fmt.Sprintf("The following query parameters are not recognised: %v", strings.Join(invalidKeys, ","))}
	}

	var result json.RawMessage
	var err error
	slog.Info("GetSessions", "date", c.QueryParam("date"), "backline", c.QueryParam("backline"), "genre", c.QueryParam("genre"))
	if c.QueryParam("date") == "" && c.QueryParam("backline") == "" && c.QueryParam("genre") == "" { // no params used, no filter
		result, err = queries.GetAllSessionsAsGeoJSON(ctx)
	} else if c.QueryParam("date") != "" && c.QueryParam("backline") != "" && c.QueryParam("genre") != "" { // all available query params are used
		result, err = queries.GetSessionsByDateAndGenreAndBacklineAsGeoJSON(ctx, dbutils.GetSessionsByDateAndGenreAndBacklineAsGeoJSONParams{
			Genres:   []string{genre},
			Backline: backline,
			Date:     pgtype.Date{Time: date, Valid: true},
		})
	} else if c.QueryParam("date") != "" && c.QueryParam("backline") != "" && c.QueryParam("genre") == "" { // date and backline are used
		result, err = queries.GetSessionsByDateAndBacklineAsGeoJSON(ctx, dbutils.GetSessionsByDateAndBacklineAsGeoJSONParams{
			Date:     pgtype.Date{Time: date, Valid: true},
			Backline: backline,
		})
	} else if c.QueryParam("date") != "" && c.QueryParam("genre") != "" && c.QueryParam("backline") == "" { // date and genre are used
		result, err = queries.GetSessionsByDateAndGenreAsGeoJSON(ctx, dbutils.GetSessionsByDateAndGenreAsGeoJSONParams{
			Date:   pgtype.Date{Time: date, Valid: true},
			Genres: []string{genre},
		})
	} else if c.QueryParam("genre") != "" && c.QueryParam("backline") != "" && c.QueryParam("date") == "" { // genre and backline are used
		result, err = queries.GetSessionsByGenreAndBacklineAsGeoJSON(ctx, dbutils.GetSessionsByGenreAndBacklineAsGeoJSONParams{
			Genres:   []string{genre},
			Backline: backline,
		})
	} else if c.QueryParam("date") != "" {
		result, err = queries.GetSessionsByDateAsGeoJSON(ctx, pgtype.Date{Time: date, Valid: true})
	} else if c.QueryParam("backline") != "" {
		result, err = queries.GetSessionsByBacklineAsGeoJSON(ctx, backline)
	} else { // genre
		result, err = queries.GetSessionsByGenreAsGeoJSON(ctx, []string{genre})
	}
	if err != nil {
		return geojson, err
	}
	json.Unmarshal(result, &geojson)
	slog.Info("GetSessions", "result", geojson)
	return geojson, nil
}

func GetSessionById(c *fuego.ContextNoBody) (types.SessionFeature[types.SessionPropertiesWithVenue], error) {
	var geojson types.SessionFeature[types.SessionPropertiesWithVenue]
	slog.Info("GetSessionById", "id", c.PathParam("id"))
	id, err := strconv.Atoi(c.PathParam("id"))
	if err != nil {
		return types.SessionFeature[types.SessionPropertiesWithVenue]{}, fuego.BadRequestError{Detail: fmt.Sprintf("Please provide a numeric ID ('/jamsession/{id}'), got: %v", c.PathParam("id"))}
	}
	result, err := queries.GetSessionByIdAsGeoJSON(ctx, int32(id))
	if err != nil {
		return types.SessionFeature[types.SessionPropertiesWithVenue]{}, err
	}
	err = json.Unmarshal([]byte(result.(string)), &geojson)
	if err != nil {
		return types.SessionFeature[types.SessionPropertiesWithVenue]{}, err
	}
	slog.Info("GetSessionById", "result", geojson)
	return geojson, nil
}

func GetVenueById(c *fuego.ContextNoBody) (types.VenueFeature, error) {
	var geojson types.VenueFeature
	slog.Info("GetVenueById", "id", c.PathParam("id"))
	id, err := strconv.Atoi(c.PathParam("id"))
	if err != nil {
		return types.VenueFeature{}, fuego.BadRequestError{Detail: fmt.Sprintf("Please provide a numeric ID ('/venues/{id}'), got: %v", c.PathParam("id"))}
	}
	result, err := queries.GetVenueByIdAsGeoJSON(ctx, int32(id))
	if err != nil {
		return types.VenueFeature{}, err
	}
	err = json.Unmarshal([]byte(result.(string)), &geojson)
	if err != nil {
		return types.VenueFeature{}, err
	}
	slog.Info("GetVenueById", "result", geojson)
	return geojson, nil
}

// the following handlers don't directly apply changes but rather prepare commits for the admin to manually run (make migrations or scripts/run-migrations.sh)
// this is to prevent users from directly modifying the database

// helper - https://github.com/golang/go/issues/63309
func ptr[T any](t T) *T { return &t }

func PostSession(c *fuego.ContextWithBody[types.SessionPropertiesWithVenue]) (types.SessionFeature[types.SessionProperties], error) {
	payload, err := c.Body()
	slog.Info("PostSession", "payload", payload)
	if err != nil {
		slog.Error("PostSession", "msg", err)
		return types.SessionFeature[types.SessionProperties]{}, errors.New("an unknown error occured")
	}
	var cmd string
	var title string
	var sessionJson []byte
	if payload.VenueName != nil { // if venue fields are present in the payload, we create a new venue in the same transaction
		venueJson, err := json.Marshal(payload.VenueProperties)
		if err != nil {
			slog.Error("PostSession", "msg", err, "props", "venue")
			return types.SessionFeature[types.SessionProperties]{}, errors.New("an unknown error occured")
		}

		// set venue to bash variable that will be evaluated as the real venue id during runtime
		payload.Venue = ptr(int32(999999))
		sessionJson, err = json.Marshal(payload.SessionProperties)
		if err != nil {
			slog.Error("PostSession", "msg", err, "props", "session")
			return types.SessionFeature[types.SessionProperties]{}, errors.New("an unknown error occured")
		}
		sessionJson = []byte(strings.Replace(string(sessionJson), "999999", "$new_id", -1))

		cmd = fmt.Sprintf("new_id=$(dbcli insert venue '%s');\n"+`dbcli insert session "%s";`, strings.ReplaceAll(string(venueJson), "'", `\'`), strings.ReplaceAll(string(sessionJson), `"`, `\"`))
		title = fmt.Sprintf("insert_venue_%v_session_%v", *payload.VenueName, *payload.SessionName)
		slog.Info("PostSession", "mode", "sessionAndVenue", "cmd", cmd)
	} else {
		sessionJson, err = json.Marshal(payload.SessionProperties)
		if err != nil {
			slog.Error("PostSession", "msg", err, "props", "session")
			return types.SessionFeature[types.SessionProperties]{}, errors.New("an unknown error occured")
		}
		cmd = fmt.Sprintf("dbcli insert session '%v'", strings.ReplaceAll(string(sessionJson), "'", `\'`))
		title = fmt.Sprintf("insert_session_%v", *payload.SessionName)
		slog.Info("PostSession", "mode", "sessionOnly", "cmd", cmd)
	}
	if err := writeMigration(cmd, title, migrationsDirectory); err != nil {
		return types.SessionFeature[types.SessionProperties]{}, err
	}
	c.SetStatus(201)
	return types.SessionFeature[types.SessionProperties]{}, nil
}

func PatchSessionById(c *fuego.ContextWithBody[types.SessionProperties]) (types.SessionFeature[types.SessionProperties], error) {
	slog.Info("PatchSessionById", "id", c.PathParam("id"))
	id, err := strconv.Atoi(c.PathParam("id"))
	if err != nil {
		return types.SessionFeature[types.SessionProperties]{}, fuego.BadRequestError{Detail: fmt.Sprintf("Please provide a numeric ID ('/jamsession/{id}'), got: %v", c.PathParam("id"))}
	}
	payload, err := c.Body()
	if err != nil {
		slog.Error("PatchSessionById", "id", id, "msg", err)
		return types.SessionFeature[types.SessionProperties]{}, errors.New("an unknown error occured")
	}
	j, err := json.Marshal(payload)
	if err != nil {
		slog.Error("PatchSessionById", "id", id, "msg", err)
		return types.SessionFeature[types.SessionProperties]{}, errors.New("an unknown error occured")
	}
	cmd := fmt.Sprintf("dbcli update session %v '%s'", id, strings.ReplaceAll(string(j), "'", `\'`))
	return types.SessionFeature[types.SessionProperties]{}, writeMigration(cmd, fmt.Sprintf("update_session_%v", id), migrationsDirectory)
}

type CommentBody struct {
	Session *int   `json:"session"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

func PostCommentForSessionById(c *fuego.ContextWithBody[CommentBody]) (types.SessionFeature[types.SessionProperties], error) {
	slog.Info("PostCommentForSessionById", "id", c.PathParam("id"))
	id, err := strconv.Atoi(c.PathParam("id"))
	if err != nil {
		return types.SessionFeature[types.SessionProperties]{}, fuego.BadRequestError{Detail: fmt.Sprintf("Please provide a numeric ID ('/jamsession/{id}'), got: %v", c.PathParam("id"))}
	}
	payload, err := c.Body()
	if err != nil {
		return types.SessionFeature[types.SessionProperties]{}, err
	}
	payload.Session = ptr(id)
	j, err := json.Marshal(payload)
	if err != nil {
		return types.SessionFeature[types.SessionProperties]{}, err
	}
	cmd := fmt.Sprintf("dbcli insert comment '%s'", strings.ReplaceAll(string(j), "'", `\'`))
	if err := writeMigration(cmd, fmt.Sprintf("insert_comment_session_%v", id), migrationsDirectory); err != nil {
		return types.SessionFeature[types.SessionProperties]{}, err
	}
	c.SetStatus(201)
	return types.SessionFeature[types.SessionProperties]{}, nil
}

func PostSuggestionsForSessionById(c *fuego.ContextWithBody[CommentBody]) (types.SessionFeature[types.SessionProperties], error) {
	slog.Info("PostSuggestionsSessionById", "id", c.PathParam("id"))
	id, err := strconv.Atoi(c.PathParam("id"))
	if err != nil {
		return types.SessionFeature[types.SessionProperties]{}, fuego.BadRequestError{Detail: fmt.Sprintf("Please provide a numeric ID ('/jamsession/{id}'), got: %v", c.PathParam("id"))}
	}

	body, err := c.Body()
	if err != nil {
		return types.SessionFeature[types.SessionProperties]{}, err
	}

	fp := filepath.Join(suggestionsDirectory, fmt.Sprintf("%v_session_%v", time.Now().Format(time.RFC3339), (id)))
	slog.Info("writing suggestion", "filepath", fp)
	os.WriteFile(fp, []byte(fmt.Sprintf("Session %v: %v", id, body)), fs.FileMode(int(0755)))
	return types.SessionFeature[types.SessionProperties]{}, nil
}

func GetCommentsBySessionId(c *fuego.ContextNoBody) ([]dbutils.GetCommentsBySessionIdRow, error) {
	slog.Info("GetCommentsBySessionId", "id", c.PathParam("id"))
	id, err := strconv.Atoi(c.PathParam("id"))
	if err != nil {
		return []dbutils.GetCommentsBySessionIdRow{}, fuego.BadRequestError{Detail: fmt.Sprintf("Please provide a numeric ID ('/jamsession/{id}/comments'), got: %v", c.PathParam("id"))}
	}
	res, err := queries.GetCommentsBySessionId(ctx, int32(id))
	if err != nil {
		slog.Error("GetCommentsBySessionId", "id", id, "err", err)
		return []dbutils.GetCommentsBySessionIdRow{}, errors.New("an unknown error occured")
	}
	return res, nil
}

func DeleteSessionById(c *fuego.ContextNoBody) (types.SessionFeature[types.SessionProperties], error) {
	slog.Info("DeleteSessionById", "id", c.PathParam("id"))

	id, err := strconv.Atoi(c.PathParam("id"))
	if err != nil {
		return types.SessionFeature[types.SessionProperties]{}, fuego.BadRequestError{Detail: fmt.Sprintf("Please provide a numeric ID ('/jamsession/{id}'), got: %v", c.PathParam("id"))}
	}
	cmd := fmt.Sprintf("dbcli delete session %v", id)
	return types.SessionFeature[types.SessionProperties]{}, writeMigration(cmd, fmt.Sprintf("delete_session_%v", id), migrationsDirectory)
}

func PostVenue(c *fuego.ContextWithBody[types.VenueProperties]) (types.VenueFeature, error) {
	payload, err := c.Body()
	if err != nil {
		slog.Error("PostVenue", "msg", err)
		return types.VenueFeature{}, errors.New("an unknown error occured")
	}
	j, err := json.Marshal(payload)
	if err != nil {
		slog.Error("PostVenue", "msg", err)
		return types.VenueFeature{}, errors.New("an unknown error occured")
	}
	cmd := fmt.Sprintf("dbcli insert venue '%s'", strings.ReplaceAll(string(j), "'", `\'`))
	if err := writeMigration(cmd, "insert_venue_"+*payload.VenueName, migrationsDirectory); err != nil {
		return types.VenueFeature{}, err
	}
	c.SetStatus(201)
	return types.VenueFeature{}, nil

}

func PatchVenueById(c *fuego.ContextWithBody[types.VenueProperties]) (types.VenueFeature, error) {
	slog.Info("PatchVenueById", "id", c.PathParam("id"))

	id, err := strconv.Atoi(c.PathParam("id"))
	if err != nil {
		return types.VenueFeature{}, fuego.BadRequestError{Detail: fmt.Sprintf("Please provide a numeric ID ('/jamsession/{id}'), got: %v", c.PathParam("id"))}
	}
	payload, err := c.Body()
	if err != nil {
		slog.Error("PatchVenueById", "id", id, "msg", err)
		return types.VenueFeature{}, errors.New("an unknown error occured")
	}
	j, err := json.Marshal(payload)
	if err != nil {
		slog.Error("PatchVenueById", "id", id, "msg", err)
		return types.VenueFeature{}, errors.New("an unknown error occured")
	}
	cmd := fmt.Sprintf("dbcli update venue %v '%s'", id, strings.ReplaceAll(string(j), "'", `\'`))
	return types.VenueFeature{}, writeMigration(cmd, fmt.Sprintf("update_venue_%v", id), migrationsDirectory)
}

func DeleteVenueById(c *fuego.ContextNoBody) (types.VenueFeature, error) {
	slog.Info("DeleteVenueById", "id", c.PathParam("id"))

	id, err := strconv.Atoi(c.PathParam("id"))
	if err != nil {
		return types.VenueFeature{}, fuego.BadRequestError{Detail: fmt.Sprintf("Please provide a numeric ID ('/jamsession/{id}'), got: %v", c.PathParam("id"))}
	}
	cmd := fmt.Sprintf("dbcli delete venue %v", id)
	return types.VenueFeature{}, writeMigration(cmd, fmt.Sprintf("delete_venue_%v", id), migrationsDirectory)
}

// helper func - writes the provided command to a timestamped file in the $MIGRATIONS_DIRECTORY. 'route' is solely for logging purposes.
func writeMigration(cmd string, title string, migrationsDirectory string) error {
	if migrationsDirectory == "" {
		slog.Error("migrationsDirectory is not set")
		return errors.New("an unknown error occured")
	}

	fp := filepath.Join(migrationsDirectory, fmt.Sprintf("%v_%v_%v.sh", time.Now().UTC().Format("20060102_150405"), time.Now().Nanosecond(), strings.ReplaceAll(strings.ReplaceAll(title, " ", "_"), "'", "")))
	slog.Info("writing migration", "filepath", fp)
	os.WriteFile(fp, []byte("#!/usr/bin/env bash\n\n"+cmd), fs.FileMode(int(0755)))
	return nil
}
