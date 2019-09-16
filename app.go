package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-redis/redis"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// global config struct holding database connection info
type config struct {
	sqlUsername   string
	sqlPassword   string
	sqlHost       string
	sqlPort       string
	sqlDbName     string
	redisPassword string
	redisHost     string
	redisPort     string
	redisDb       int
}

// method to initialize config struct from environment variables
func (conf *config) configure() {
	conf.sqlUsername = os.Getenv("THORCAST_DB_USERNAME")
	conf.sqlPassword = os.Getenv("THORCAST_DB_PASSWORD")
	conf.sqlHost = os.Getenv("THORCAST_DB_HOST")
	conf.sqlPort = os.Getenv("THORCAST_DB_PORT")
	conf.sqlDbName = os.Getenv("THORCAST_DB_NAME")
	conf.redisPassword = os.Getenv("REDIS_PASSWORD")
	conf.redisHost = os.Getenv("REDIS_HOST")
	conf.redisPort = os.Getenv("REDIS_PORT")
	conf.redisDb, _ = strconv.Atoi(os.Getenv("REDIS_DB"))
}

var conf = config{}

type App struct {
	Router *mux.Router
	Logger http.Handler
	DB     *sql.DB
	Redis  *redis.Client
}

func (a *App) HourlyForecast(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	checkHours, ok := vars["hours"]
	if !ok {
		checkHours = "12"
	}
	city, state, hours, err := SanitizeHourlyInputs(
		vars["city"],
		vars["state"],
		checkHours,
	)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		l := Location{City: city.asName, State: state.asName}
		var hourlyForecasts []string
		hourlyForecasts, err := a.LookupHourlyForecast(city, state, hours)
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
					l.SetCoords(FetchCoords(city.asURL, state.asURL))
					if err = a.RegisterLocation(l); err != nil {
						a.IncrementLocation(l)
					}
				default:
					respondWithError(w, http.StatusBadRequest, err.Error())
				}
			} else {
				a.IncrementLocation(l)
			}
			forecastURL := FetchHourlyForecastURL(l)
			log.Printf("%s", forecastURL)
			forecasts := FetchForecasts(forecastURL)
			hourlyForecasts = a.CacheHourlyForecasts(city, state, hours, forecasts)
		} else if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		} else {
			a.IncrementLocation(l)
		}
		resp := map[string]string{
			"forecast": strings.Join(hourlyForecasts, "\n"),
			"city":     city.asName,
			"state":    state.asName,
			"hours":    checkHours}
		respondWithJSON(w, http.StatusOK, resp)
	}
}

// Forecast returns the detailed forecast for a given city, state, and period
// if period is not specified in the HTTP request, it defaults to today
func (a *App) DetailedForecast(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	checkPeriod, ok := vars["period"]
	if !ok {
		checkPeriod = "today"
	}

	city, state, period, err := SanitizeDetailedInputs(
		vars["city"],
		vars["state"],
		checkPeriod,
	)
	// l := Location{City: city.asName, State: state.asName}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		l := Location{City: city.asName, State: state.asName}
		var forecast string
		forecast, err := a.LookupDetailedForecast(city, state, period)
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
					l.SetCoords(FetchCoords(city.asURL, state.asURL))
					if err = a.RegisterLocation(l); err != nil {
						a.IncrementLocation(l)
					}
				default:
					respondWithError(w, http.StatusBadRequest, err.Error())
				}
			} else {
				a.IncrementLocation(l)
			}
			forecastURL := FetchDetailedForecastURL(l)
			forecasts := FetchForecasts(forecastURL)
			forecast = a.CacheDetailedForecasts(city, state, period, forecasts)
		} else if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		} else {
			a.IncrementLocation(l)
		}
		resp := map[string]string{
			"forecast": forecast,
			"city":     city.asName,
			"state":    state.asName,
			"period":   period.asName}
		respondWithJSON(w, http.StatusOK, resp)
	}
}

// Random provides a forecast for a random city, state, and period
// city and state are determined by selecting a random location from the database
// period is selected randomly within the next week
func (a *App) RandomDetailedForecast(w http.ResponseWriter, r *http.Request) {
	var l Location
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
		log.Fatal(err)
	}

	period := randomPeriod()
	city := sanitizeCity(l.City)
	state, _ := sanitizeState(l.State)
	forecast, err := a.LookupDetailedForecast(city, state, period)
	if err == redis.Nil {
		forecastURL := FetchDetailedForecastURL(l)
		forecasts := FetchForecasts(forecastURL)
		forecast = a.CacheDetailedForecasts(city, state, period, forecasts)
	}
	a.IncrementLocation(l)
	resp := map[string]string{
		"forecast": forecast,
		"city":     city.asName,
		"state":    state.asName,
		"period":   period.asName}
	respondWithJSON(w, http.StatusOK, resp)
}

// Initializes all endpoints for the api
func (a *App) InitializeRoutes() {
	a.Router.HandleFunc("/api/forecast/detailed/city={city}&state={state}&period={period}", a.DetailedForecast).Methods("GET")
	a.Router.HandleFunc("/api/forecast/detailed/city={city}&state={state}", a.DetailedForecast).Methods("GET")
	a.Router.HandleFunc("/api/forecast/detailed/random", a.RandomDetailedForecast).Methods("GET")
	a.Router.HandleFunc("/api/forecast/hourly/city={city}&state={state}&hours={hours}", a.HourlyForecast).Methods("GET")
	a.Router.HandleFunc("/api/forecast/hourly/city={city}&state={state}", a.HourlyForecast).Methods("GET")
}

// Initializes the application as a whole
func (a *App) Initialize() {
	var err error
	conf.configure()
	sqlDataSource := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		conf.sqlUsername,
		conf.sqlPassword,
		conf.sqlHost,
		conf.sqlPort,
		conf.sqlDbName)
	a.DB, err = sql.Open("postgres", sqlDataSource)
	if err != nil {
		log.Fatal(err)
	}
	a.Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.redisHost, conf.redisPort),
		Password: conf.redisPassword,
		DB:       conf.redisDb,
	})
	a.Router = mux.NewRouter()
	a.Logger = handlers.CombinedLoggingHandler(os.Stdout, a.Router)
	a.InitializeRoutes()
}

// Starts the app to listen on the port specitied by the env variable SERVER_PORT
func (a *App) Run() {
	port := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	log.Fatal(http.ListenAndServe(port, a.Logger))
}
