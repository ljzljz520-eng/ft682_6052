package main

import (
	"log"
	"os"

	"venue-reservation/internal/cli"
	"venue-reservation/internal/httpapi"
	"venue-reservation/internal/seed"
	"venue-reservation/internal/service"
	"venue-reservation/internal/store"
)

func main() {
	config, err := cli.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	database, err := store.Open(config.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	app := service.New(database)
	if cli.IsSeedEnabled(config) {
		if err := seed.EnsureDemoData(app); err != nil {
			log.Fatal(err)
		}
	}
	server := httpapi.NewServer(app)
	log.Printf("listening on %s", config.Address)
	if err := server.ListenAndServe(config.Address); err != nil {
		log.Fatal(err)
	}
}
