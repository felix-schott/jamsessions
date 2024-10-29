package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	dbutils "github.com/felix-schott/jamsessions/backend/internal/db"
	"github.com/felix-schott/jamsessions/backend/internal/geocoding"

	"github.com/alexflint/go-arg"
)

type UpdateCmd struct {
	Table   string `arg:"positional"`
	Id      int    `arg:"positional"`
	Payload string `arg:"positional"`
}
type InsertCmd struct {
	Table   string `arg:"positional"`
	Payload string `arg:"positional"`
}
type DeleteCmd struct {
	Table string `arg:"positional"`
	Id    int    `arg:"positional"`
}

type args struct {
	Update *UpdateCmd `arg:"subcommand:update"`
	Insert *InsertCmd `arg:"subcommand:insert"`
	Delete *DeleteCmd `arg:"subcommand:delete"`
}

func (args) Description() string {
	return "CLI tool to modify the session/venue database."
}

var ctx = context.Background()
var queries *dbutils.Queries

func main() {

	pool, err := dbutils.CreatePool(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	queries = dbutils.New(pool)

	var args args
	p := arg.MustParse(&args)

	switch {
	case args.Update != nil:
		switch args.Update.Table {
		case "venue":
			log.Printf("Updating record %v of table venue\n", args.Update.Id)
			var payload dbutils.UpdateVenueByIdParams
			if err := json.Unmarshal([]byte(args.Update.Payload), &payload); err != nil {
				p.Fail(fmt.Sprintf("couldn't parse payload with the following error: %v", err))
			}
			payload.VenueID = int32(args.Update.Id)

			// geocode address
			if payload.AddressFirstLine != nil || payload.AddressSecondLine != nil || payload.City != nil || payload.Postcode != nil {
				log.Println("Obtaining coordinates from address")
				var s string
				if payload.AddressSecondLine != nil {
					s = *payload.AddressFirstLine + *payload.AddressSecondLine
				} else {
					s = *payload.AddressFirstLine
				}
				loc, err := geocoding.Geocode(s, *payload.City, *payload.Postcode)
				if err != nil {
					log.Fatalf("failed to obtain coordinates from provided address: %v", err)
				}
				payload.Geom = loc
			}
			if err := queries.UpdateVenueById(ctx, payload); err != nil {
				log.Fatalf("failed to run query: %v", err)
			}
		case "session":
			fmt.Printf("Updating record %v of table session\n", args.Update.Id)
			var payload dbutils.UpdateJamSessionByIdParams
			if err := json.Unmarshal([]byte(args.Update.Payload), &payload); err != nil {
				p.Fail(fmt.Sprintf("couldn't parse payload with the following error: %v", err))
			}
			payload.SessionID = int32(args.Update.Id)
			if err := queries.UpdateJamSessionById(ctx, payload); err != nil {
				log.Fatalf("failed to run query: %v", err)
			}
		default:
			p.Fail(fmt.Sprintf("available tables: 'venue' or 'session', got %v", args.Update.Table))
		}
	case args.Insert != nil:
		switch args.Insert.Table {
		case "venue":
			log.Println("Inserting record into table venue")
			var payload dbutils.InsertVenueParams
			if err := json.Unmarshal([]byte(args.Insert.Payload), &payload); err != nil {
				p.Fail(fmt.Sprintf("couldn't parse payload with the following error: %v", err))
			}

			// geocode address
			if payload.AddressFirstLine != "" || payload.AddressSecondLine != nil || payload.City != "" || payload.Postcode != "" {
				log.Println("Obtaining coordinates from address")
				var s string
				if payload.AddressSecondLine != nil {
					s = payload.AddressFirstLine + *payload.AddressSecondLine
				} else {
					s = payload.AddressFirstLine
				}
				loc, err := geocoding.Geocode(s, payload.City, payload.Postcode)
				if err != nil {
					log.Fatalf("failed to obtain coordinates from provided address: %v", err)
				}
				payload.Geom = loc
			}

			newId, err := queries.InsertVenue(ctx, payload)
			if err != nil {
				log.Fatalf("failed to run query: %v", err)
			}
			log.Printf("Inserted record with ID %v into table venue\n", newId)
			fmt.Print(newId) // write new id to stdout
		case "session":
			log.Println("Inserting record into table session")
			var payload dbutils.InsertJamSessionParams
			if err := json.Unmarshal([]byte(args.Insert.Payload), &payload); err != nil {
				p.Fail(fmt.Sprintf("couldn't parse payload with the following error: %v", err))
			}
			newId, err := queries.InsertJamSession(ctx, payload)
			if err != nil {
				log.Fatalf("failed to run query: %v", err)
			}
			log.Printf("Inserted record with ID %v into table session\n", newId)
			fmt.Print(newId) // write new id to stdout
		case "comment":
			log.Println("Inserting record into table comments")
			var payload dbutils.InsertSessionCommentParams
			if err := json.Unmarshal([]byte(args.Insert.Payload), &payload); err != nil {
				p.Fail(fmt.Sprintf("couldn't parse payload with the following error: %v", err))
			}
			newId, err := queries.InsertSessionComment(ctx, payload)
			if err != nil {
				log.Fatalf("failed to run query: %v", err)
			}
			log.Printf("Inserted record with ID %v into table comments\n", newId)
			fmt.Print(newId) // write new id to stdout
		case "rating":
			log.Println("Inserting record into table ratings")
			var payload dbutils.InsertSessionRatingParams
			if err := json.Unmarshal([]byte(args.Insert.Payload), &payload); err != nil {
				p.Fail(fmt.Sprintf("couldn't parse payload with the following error: %v", err))
			}
			newId, err := queries.InsertSessionRating(ctx, payload)
			if err != nil {
				log.Fatalf("failed to run query: %v", err)
			}
			log.Printf("Inserted record with ID %v into table ratings\n", newId)
			fmt.Print(newId) // write new id to stdout
		default:
			p.Fail(fmt.Sprintf("available tables: 'venue' or 'session', got %v", args.Insert.Table))
		}
	case args.Delete != nil:
		switch args.Delete.Table {
		case "venue":
			log.Printf("Deleting record %v from table venue\n", args.Delete.Id)
			if err := queries.DeleteVenueById(ctx, int32(args.Delete.Id)); err != nil {
				log.Fatalf("failed to run query: %v", err)
			}
		case "session":
			log.Printf("Deleting record %v from table session\n", args.Delete.Id)
			if err := queries.DeleteJamSessionById(ctx, int32(args.Delete.Id)); err != nil {
				log.Fatalf("failed to run query: %v", err)
			}
		default:
			p.Fail(fmt.Sprintf("available tables: 'venue' or 'session', got %v", args.Delete.Table))
		}
	}
}
