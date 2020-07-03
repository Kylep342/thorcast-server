package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

const gcAPI = "https://maps.googleapis.com/maps/api/geocode/json"

// Struct holding data from the maps.google.com geocoding api
type geocodeAPIResp struct {
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport     struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		PartialMatch bool     `json:"partial_match"`
		PlaceID      string   `json:"place_id"`
		Types        []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}

// Struct holding the lat, lng pair from a maps.google.com geocode api response
type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// FetchCoords returns coordinates for a given address
func FetchCoords(city, state string) (Coordinates, error) {
	APIKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	requestURL := fmt.Sprintf("%s?address=%s,%s&key=%s", gcAPI, city, state, APIKey)
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Google Maps API error is: %s\n", err.Error())
	}
	var geocode geocodeAPIResp
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&geocode)
	if err != nil {
		log.Printf("JSON decoding error is: %s\n", err.Error())
	}
	log.Printf("Status is %s\n", geocode.Status)
	if geocode.Status != "OK" {
		switch geocode.Status {
		case "ZERO_RESULTS":
			return Coordinates{}, errors.New("location not found")
		default:
			return Coordinates{}, errors.New("internal error")
		}
	}
	if geocode.Results[0].PartialMatch {
		return Coordinates{}, errors.New("location not found")
	}
	coords := Coordinates{Lat: geocode.Results[0].Geometry.Location.Lat, Lng: geocode.Results[0].Geometry.Location.Lng}
	return coords, nil
}
