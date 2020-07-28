package app

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/go-redis/redis"
	"github.com/kylep342/thorcast-server/cache"
	"github.com/kylep342/thorcast-server/db"
	"github.com/kylep342/thorcast-server/utils"
)

// Custom404Handler defines a catchall response for invalid API endpoints
func (a *App) Custom404Handler(w http.ResponseWriter, r *http.Request) {
	code := http.StatusNotFound
	respondWithError(w, code, http.StatusText(code))
}

// HourlyForecast returns hourly forecast data for the specified city, state, and duration
func (a *App) HourlyForecastHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	var checkHours string

	hasHours, ok := params["hours"]
	if !ok {
		checkHours = "12"
	} else {
		checkHours = hasHours[0]
	}

	city, state, hours, err := utils.SanitizeHourlyInputs(
		params["city"][0],
		params["state"][0],
		checkHours,
	)
	if err != nil {
		log.Printf("Error checking client inputs: %s\n", err.Error())
		code := http.StatusBadRequest
		respondWithError(w, code, http.StatusText(code))
	} else {
		l := models.Location{City: city.asName, State: state.asName}
		var hourlyForecasts []string
		hourlyForecasts, err := cache.LookupHourlyForecast(a.Redis, city, state, hours)
		if err == redis.Nil {
			row := a.DB.QueryRow(
				`SELECT
					lat,
					lng
				FROM geocodex
				WHERE LOWER(city) = LOWER($1)
				AND state = $2
				;`,
				l.City,
				l.State)
			if err := row.Scan(&l.Lat, &l.Lng); err != nil {
				switch err {
				case sql.ErrNoRows:
					coords, err := FetchCoords(city.asURL, state.asURL)
					if err != nil {
						code := http.StatusNotFound
						respondWithError(w, code, http.StatusText(code))
					} else {
						l.SetLocationCoordinates(coords)
						if err = db.RegisterLocation(a.DB, l); err != nil {
							// a.IncrementLocation(l)
							code := http.StatusInternalServerError
							respondWithError(w, code, http.StatusText(code))
						}
					}
				default:
					log.Printf("Error scanning lat/lng from the database: %s\n", err.Error())
					code := http.StatusInternalServerError
					respondWithError(w, code, http.StatusText(code))
				}
			} else {
				db.IncrementLocation(a.DB, l)
			}
			forecastURL, err := apis.FetchHourlyForecastURL(l)
			if err != nil {
				code := http.StatusInternalServerError
				respondWithError(w, code, http.StatusText(code))
			}
			forecasts, err := apis.FetchForecasts(forecastURL)
			if err != nil {
				log.Printf("Error when fetching forecasts\nError is %s\n", err.Error())
				code := http.StatusInternalServerError
				respondWithError(w, code, http.StatusText(code))
			}
			hourlyForecasts = cache.CacheHourlyForecasts(a.Redis, city, state, hours, forecasts)
		} else if err != nil {
			log.Printf("Error looking up hourly forecasts: %s\n", err.Error())
			code := http.StatusInternalServerError
			respondWithError(w, code, http.StatusText(code))
		} else {
			db.IncrementLocation(a.DB, l)
		}
		resp := map[string]string{
			"forecast": strings.Join(hourlyForecasts, "\n"),
			"city":     city.asName,
			"state":    state.asName,
			"hours":    checkHours}
		responses.RespondWithJSON(w, http.StatusOK, resp)
	}
}

// DetailedForecast returns the detailed forecast for a given city, state, and period
// if period is not specified in the HTTP request, it defaults to today
func (a *App) DetailedForecastHandler(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()

	var checkPeriod string

	hasPeriod, ok := params["period"]
	if !ok {
		checkPeriod = "today"
	} else {
		checkPeriod = hasPeriod[0]
	}

	city, state, period, err := utils.SanitizeDetailedInputs(
		params["city"][0],
		params["state"][0],
		checkPeriod,
	)
	if err != nil {
		log.Printf("Error sanitizing client inputs: %s\n", err.Error())
		code := http.StatusBadRequest
		respondWithError(w, code, http.StatusText(code))
	} else {
		l := models.Location{City: city.asName, State: state.asName}
		var forecast string
		forecast, err := cache.LookupDetailedForecast(a.Redis, city, state, period)
		if err == redis.Nil {
			row := a.DB.QueryRow(
				`SELECT
				lat,
				lng
			FROM geocodex
			WHERE LOWER(city) = LOWER($1)
			AND state = $2
			;`,
				l.City,
				l.State)
			if err := row.Scan(&l.Lat, &l.Lng); err != nil {
				log.Printf("row scan error: %s\n", err.Error())
				switch err {
				case sql.ErrNoRows:
					log.Printf("city: %s state: %s\n", city.asURL, state.asURL)
					coords, err := FetchCoords(city.asURL, state.asURL)
					if err != nil {
						code := http.StatusNotFound
						respondWithError(w, code, http.StatusText(code))
					} else {
						l.SetLocationCoordinates(coords)
						if err = db.RegisterLocation(a.DB, l); err != nil {
							// a.IncrementLocation(l)
							code := http.StatusInternalServerError
							respondWithError(w, code, http.StatusText(code))
						}
					}
				default:
					log.Printf("error is: %s\n", err.Error())
					code := http.StatusBadRequest
					respondWithError(w, code, http.StatusText(code))
				}
			} else {
				db.IncrementLocation(a.DB, l)
			}
			forecastURL, err := apis.FetchDetailedForecastURL(l)
			if err != nil {
				code := http.StatusInternalServerError
				respondWithError(w, code, http.StatusText(code))
			}
			forecasts, err := apis.FetchForecasts(forecastURL)
			if err != nil {
				log.Printf("Error when fetching forecasts\nError is %s\n", err.Error())
				code := http.StatusInternalServerError
				respondWithError(w, code, http.StatusText(code))
			}
			forecast = cache.CacheDetailedForecasts(a.Redis, city, state, period, forecasts)
		} else if err != nil {
			log.Printf("Error looking up detailed forecast: %s\n", err.Error())
			code := http.StatusInternalServerError
			respondWithError(w, code, http.StatusText(code))
		} else {
			db.IncrementLocation(a.DB, l)
		}
		resp := map[string]string{
			"forecast": forecast,
			"city":     city.asName,
			"state":    state.asName,
			"period":   period.asName}
		responses.RespondWithJSON(w, http.StatusOK, resp)
	}
}

// RandomDetailedForecast provides a forecast for a random city, state, and period
// city and state are determined by selecting a random location from the database
// period is selected randomly within the next week
func (a *App) RandomDetailedForecastHandler(w http.ResponseWriter, r *http.Request) {
	var l models.Location
	var forecast string
	row := a.DB.QueryRow(
		`SELECT
			city,
			state,
			lat,
			lng
		FROM geocodex
		ORDER BY random()
		LIMIT 1;`)
	if err := row.Scan(&l.City, &l.State, &l.Lat, &l.Lng); err != nil {
		log.Printf("Error when reading geocodex information from the database.\nError is %s\n", err.Error())
		code := http.StatusInternalServerError
		respondWithError(w, code, http.StatusText(code))
	}

	period := utils.RandomPeriod()
	city := utils.SanitizeCity(l.City)
	state, _ := utils.SanitizeState(l.State)
	forecast, err := cache.LookupDetailedForecast(a.Redis, city, state, period)
	if err == redis.Nil {
		forecastURL, err := apis.FetchDetailedForecastURL(l)
		if err != nil {
			code := http.StatusInternalServerError
			respondWithError(w, code, http.StatusText(code))
		}
		forecasts, err := apis.FetchForecasts(forecastURL)
		if err != nil {
			log.Printf("Error when fetching forecasts\nError is %s\n", err.Error())
			code := http.StatusInternalServerError
			respondWithError(w, code, http.StatusText(code))
		}
		forecast = cache.CacheDetailedForecasts(a.Redis, city, state, period, forecasts)
	}
	db.IncrementLocation(a.DB, l)
	resp := map[string]string{
		"forecast": forecast,
		"city":     city.asName,
		"state":    state.asName,
		"period":   period.asName}
	responses.RespondWithJSON(w, http.StatusOK, resp)
}
