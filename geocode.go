package main

import (
	"encoding/json"
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
		PlaceID string   `json:"place_id"`
		Types   []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}

// Struct holding the lat, lng pair from a maps.google.com geocode api response
type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// FetchCoords returns coordinates for a given address
func FetchCoords(city, state string) Coordinates {
	APIKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	requestURL := fmt.Sprintf("%s?address=%s,%s&key=%s", gcAPI, city, state, APIKey)
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Fatal(err)
	}
	var geocode geocodeAPIResp
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&geocode)
	if err != nil {
		log.Fatal(err)
	}
	coord := Coordinates{Lat: geocode.Results[0].Geometry.Location.Lat, Lng: geocode.Results[0].Geometry.Location.Lng}
	return coord
}
