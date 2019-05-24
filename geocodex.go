package main

import (
	"log"
	"strings"
)

func (a *App) RegisterLocation(l Location) error {
	instertStmt := `
	INSERT INTO geocodex (city, state, lat, lng, requests)
	VALUES ($1, $2, $3, $4, 1)`
	_, err := a.DB.Exec(instertStmt, l.City, l.State, l.Lat, l.Lng)
	if err != nil {
		// should handle a race condition where request to register was valid
		// but another user concurrently made the request
		// in that case, increment the `db:"requests"` field for that location
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return err
		} else {
			log.Fatal(err)
		}
	}
	return nil
}

func (a *App) IncrementLocation(l Location) {
	updateStmt := `
	UPDATE geocodex
	SET requests = requests+1
	WHERE lower(city) = lower($1)
	AND state = $2`
	foo, err := a.DB.Exec(updateStmt, l.City, l.State)
	log.Printf("%s", foo)
	if err != nil {
		log.Fatal(err)
	}
}
