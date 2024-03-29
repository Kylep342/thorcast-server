package db

import (
	"database/sql"
	"log"

	"github.com/kylep342/thorcast-server/pkg/models"
)

// RegisterLocation persists a city, state, lat, lng group in the database
func RegisterLocation(db, *sql.DB, l models.Location) error {
	instertStmt := `
	INSERT INTO geocodex (city, state, lat, lng, requests)
	VALUES ($1, $2, $3, $4, 1)
	ON CONFLICT ON CONSTRAINT geocodex_pkey DO UPDATE
	SET requests = geocodex.requests+1
	`
	_, err := db.Exec(instertStmt, l.City, l.State, l.Lat, l.Lng)
	if err != nil {
		log.Printf("An unexpected error occurred when inserting into geocodex\nError is: %s\n", err.Error())
		return err
	}
	return nil
}

// IncrementLocation increments the requests counter of a location already stored in the database
func IncrementLocation(db *sql.DB, l models.Location) {
	updateStmt := `
	UPDATE geocodex
	SET requests = requests+1
	WHERE lower(city) = lower($1)
	AND state = $2`
	foo, err := db.Exec(updateStmt, l.City, l.State)
	log.Printf("%s", foo)
	if err != nil {
		log.Printf("An unexpected error occurred when updating requests in geocodex\nError is: %s\n", err.Error())
	}
}
