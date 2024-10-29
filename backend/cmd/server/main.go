package main

import (
	"context"
	"io/fs"
	"log"
	"os"

	dbutils "github.com/felix-schott/jamsessions/backend/internal/db"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"
	"github.com/rs/cors"
)

var queries *dbutils.Queries
var ctx = context.Background()
var suggestionsDirectory string
var migrationsDirectory string

func main() {

	// DATABASE CONNECTION
	pool, err := dbutils.CreatePool(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	queries = dbutils.New(pool)

	// cli migrations env vars
	suggestionsDirectory = os.Getenv("MIGRATIONS_SUGGESTIONS")
	if suggestionsDirectory == "" {
		log.Fatal("Please provide the environment variable MIGRATIONS_SUGGESTIONS")
	}

	if err := os.MkdirAll(suggestionsDirectory, fs.FileMode(int(0755))); err != nil {
		log.Fatalf("could not create directory %v: %v", suggestionsDirectory, err)
	}

	migrationsDirectory = os.Getenv("MIGRATIONS_DIRECTORY")
	if migrationsDirectory == "" {
		log.Fatal("Please provide the environment variable MIGRATIONS_DIRECTORY")
	}
	log.Println("The migrations directory is", migrationsDirectory)

	if err := os.MkdirAll(migrationsDirectory, fs.FileMode(int(0755))); err != nil {
		log.Fatalf("could not create directory %v: %v", migrationsDirectory, err)
	}

	// SERVER
	serverAddr := os.Getenv("SERVER_ADDRESS")
	if serverAddr == "" {
		log.Fatal("Please provide a server address (host:port) using the environment variable SERVER_ADDRESS")
	}
	s := fuego.NewServer(fuego.WithAddr(serverAddr), fuego.WithCorsMiddleware(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}).Handler))
	s.OpenApiSpec.Info = &openapi3.Info{
		Title:       "Jam Sessions",
		Description: "PLACEHOLDER",
		Version:     "0.0.1",
	}

	// MIDDLEWARE
	// fuego.Use(s, LoggingMiddleware) //, CORSMiddleWare) - for some reason the logging middleware leads to EOF errors (with post requests)

	// API VERSION 1 - Routes
	v1 := fuego.Group(s, "/v1")

	fuego.Get(v1, "/", func(c *fuego.ContextNoBody) (string, error) {
		return "Please use the versioned route /v1 (consult /swagger/index.html for interactive documentation).", nil
	})

	fuego.Get(v1, "/venues", GetVenues).Summary("Get all venues")

	fuego.Get(v1, "/venues/{id}", GetVenueById).Summary("Get a venues by its ID")

	fuego.Post(v1, "/venues", PostVenue).Summary("Add a venue")

	fuego.Patch(v1, "/venues/{id}", PatchVenueById).Summary("Update a venue by ID")

	fuego.Delete(v1, "/venues/{id}", DeleteVenueById).Summary("Delete a venue by ID")

	fuego.Get(v1, "/jamsessions", GetSessions).Summary("Get all jam sessions").Description("Use '/v1/jamsessions?date=2024-01-30' to list jamsessions happening on a particular date. The result is inferred and may not be accurate, especially for past time frames. Use '/jamsessions?backline=PA,Drums' to filter by backline provided (accepted values: 'PA', 'Drums', 'Guitar_Amp', 'Bass_Amp', 'Microphone', 'MiscPercussion'.") // TODO, explain query params

	fuego.Post(v1, "/jamsessions", PostSession).Summary("Add a jam session")

	fuego.Get(v1, "/jamsessions/{id}", GetSessionById).Summary("Get a jam session by ID")

	fuego.Patch(v1, "/jamsessions/{id}", PatchSessionById).Summary("Update a jam session by ID")

	fuego.Delete(v1, "/jamsessions/{id}", DeleteSessionById).Summary("Delete a jam session by ID")

	fuego.Post(v1, "/jamsessions/{id}/comments", PostCommentForSessionById).Summary("Post a comment for a session by ID")

	fuego.Post(v1, "/jamsessions/{id}/suggestions", PostSuggestionsForSessionById).Summary("Post feedback/suggest changes for a session by ID")

	fuego.Get(v1, "/jamsessions/{id}/comments", GetCommentsBySessionId).Summary("Get all comments for a session by ID")

	s.Run()
}
