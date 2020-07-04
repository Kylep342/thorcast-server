package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Root URL for weather.gov's api
var weatherGovAPI = os.Getenv("WEATHER_GOV_API")

// Points holds data from the request to api.weather.gov/points
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
					Value    float64 `json:"value"`
					UnitCode string  `json:"unitCode"`
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

// Forecasts holds data from the request from the Points.Properties.Forecast url
type Forecasts struct {
	Context  []interface{} `json:"@context"`
	Type     string        `json:"type"`
	Geometry struct {
		Type        string        `json:"type"`
		Coordinates [][][]float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		Updated           time.Time `json:"updated"`
		Units             string    `json:"units"`
		ForecastGenerator string    `json:"forecastGenerator"`
		GeneratedAt       time.Time `json:"generatedAt"`
		UpdateTime        time.Time `json:"updateTime"`
		ValidTimes        string    `json:"validTimes"`
		Elevation         struct {
			Value    float64 `json:"value"`
			UnitCode string  `json:"unitCode"`
		} `json:"elevation"`
		Periods []struct {
			Number           int         `json:"number"`
			Name             string      `json:"name"`
			StartTime        string      `json:"startTime"`
			EndTime          string      `json:"endTime"`
			IsDaytime        bool        `json:"isDaytime"`
			Temperature      float64     `json:"temperature"`
			TemperatureUnit  string      `json:"temperatureUnit"`
			TemperatureTrend interface{} `json:"temperatureTrend"`
			WindSpeed        string      `json:"windSpeed"`
			WindDirection    string      `json:"windDirection"`
			Icon             string      `json:"icon"`
			ShortForecast    string      `json:"shortForecast"`
			DetailedForecast string      `json:"detailedForecast"`
		} `json:"periods"`
	} `json:"properties"`
}

func fetchPoints(l Location) (Points, error) {
	requestURL := fmt.Sprintf("%s/%f,%f", weatherGovAPI, l.Lat, l.Lng)
	resp, err := http.Get(requestURL)
	log.Printf("response is %v\n", resp)
	if err != nil {
		log.Printf("Error is %s\n", err.Error())
		return Points{}, err
	}
	var p Points
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		log.Printf("Error is %s\n", err.Error())
		return Points{}, err
	}
	return p, nil
}

// FetchDetailedForecastURL extracts the URL for a forecast from the
// api.weather.gov/points response for the specified (Lat, Lng) pair
func FetchDetailedForecastURL(l Location) (string, error) {
	point, err := fetchPoints(l)
	if err != nil {
		log.Printf("Error caught.\n")
		return "", err
	}
	return point.Properties.Forecast, nil
}

// FetchHourlyForecastURL extracts the URL for an hourly forecast from the
// api.weather.gov/points response for the specified (Lat, Lng) pair
func FetchHourlyForecastURL(l Location) (string, error) {
	point, err := fetchPoints(l)
	if err != nil {
		log.Printf("Error caught.\n")
		return "", err
	}
	return point.Properties.ForecastHourly, nil
}

// FetchForecasts extract all periods of forecasts from the forecast url
// for a requested city and state
func FetchForecasts(forecastsURL string) (Forecasts, error) {
	resp, err := http.Get(forecastsURL)
	if err != nil {
		log.Printf("Error is %s\n", err.Error())
		return Forecasts{}, err
	}
	var forecasts Forecasts
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&forecasts)
	if err != nil {
		log.Printf("Error when decoding json to Forecasts.\nError is %s\n", err.Error())
		return Forecasts{}, err
	}
	return forecasts, nil
}
