package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Root URL for weather.gov's api
const weatherGovAPI = "https://api.weather.gov/points"

// Struct holding data from the request to api.weather.gov/points
type Points struct {
	Context  []interface{} `json:"@context"`
	ID       string        `json:"id"`
	Type     string        `json:"type"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		ID                  string `json:"@id"`
		Type                string `json:"@type"`
		Cwa                 string `json:"cwa"`
		ForecastOffice      string `json:"forecastOffice"`
		GridX               int    `json:"gridX"`
		GridY               int    `json:"gridY"`
		Forecast            string `json:"forecast"`
		ForecastHourly      string `json:"forecastHourly"`
		ForecastGridData    string `json:"forecastGridData"`
		ObservationStations string `json:"observationStations"`
		RelativeLocation    struct {
			Type     string `json:"type"`
			Geometry struct {
				Type        string    `json:"type"`
				Coordinates []float64 `json:"coordinates"`
			} `json:"geometry"`
			Properties struct {
				City     string `json:"city"`
				State    string `json:"state"`
				Distance struct {
					Value    float64 `json:"value"`
					UnitCode string  `json:"unitCode"`
				} `json:"distance"`
				Bearing struct {
					Value    int    `json:"value"`
					UnitCode string `json:"unitCode"`
				} `json:"bearing"`
			} `json:"properties"`
		} `json:"relativeLocation"`
		ForecastZone    string `json:"forecastZone"`
		County          string `json:"county"`
		FireWeatherZone string `json:"fireWeatherZone"`
		TimeZone        string `json:"timeZone"`
		RadarStation    string `json:"radarStation"`
	} `json:"properties"`
}

// Struct holding data from the request from the Points.Properties.Forecast url
type Forecasts struct {
	Context  []interface{} `json:"@context"`
	Type     string        `json:"type"`
	//TODO: This is not used in the app, but errors occur when parsing.
	// temporarily omitted until resolved at: https://stackoverflow.com/questions/56141009/
	//
	// Geometry struct {
	// 	Type       string `json:"type"`
	// 	Geometries []struct {
	// 		Type        string    `json:"type"`
	// 		Coordinates []float64 `json:"coordinates"`
	// 	} `json:"geometries"`
	// } `json:"geometry"`
	Properties struct {
		Updated           time.Time `json:"updated"`
		Units             string    `json:"units"`
		ForecastGenerator string    `json:"forecastGenerator"`
		GeneratedAt       time.Time `json:"generatedAt"`
		UpdateTime        time.Time `json:"updateTime"`
		ValidTimes        string 	`json:"validTimes"`
		Elevation         struct {
			Value    float64 `json:"value"`
			UnitCode string  `json:"unitCode"`
		} `json:"elevation"`
		Periods []struct {
			Number           int         	`json:"number"`
			Name             string      	`json:"name"`
			StartTime        time.Time      `json:"startTime"`
			EndTime          time.Time      `json:"endTime"`
			IsDaytime        bool        	`json:"isDaytime"`
			Temperature      int         	`json:"temperature"`
			TemperatureUnit  string      	`json:"temperatureUnit"`
			TemperatureTrend interface{} 	`json:"temperatureTrend"`
			WindSpeed        string      	`json:"windSpeed"`
			WindDirection    string      	`json:"windDirection"`
			Icon             string      	`json:"icon"`
			ShortForecast    string      	`json:"shortForecast"`
			DetailedForecast string      	`json:"detailedForecast"`
		} `json:"periods"`
	} `json:"properties"`
}

// Function to extract the URL for a forecast for the specified (Lat, Lng) pair
func FetchForecastURL(l Location) string {
	requestURL := fmt.Sprintf("%s/%f,%f", weatherGovAPI, l.Lat, l.Lng)
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Fatal(err)
	}
	var point Points
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&point)
	if err != nil {
		log.Fatal(err)
	}
	return point.Properties.Forecast
}

// Funciton to extract all periods of forecasts for a reuqested city and state
func FetchForecasts(forecastsURL string) Forecasts {
	resp, err := http.Get(forecastsURL)
	if err != nil {
		log.Fatal(err)
	}
	var forecasts Forecasts
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&forecasts)
	if err != nil {
		log.Fatal(err)
	}
	return forecasts
}

// SelectForecast extracts the desired forecast period from the response of FetchForecasts
func SelectForecast(forecasts Forecasts, period Period) string {
	for _, fc := range forecasts.Properties.Periods {
		if (fc.StartTime.Weekday().String() == period.dayOfWeek) && (fc.IsDaytime == period.isDaytime) {
			return fc.DetailedForecast
		}
	}
	return ""
}