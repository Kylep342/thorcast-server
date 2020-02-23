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

func (a *App) Custom404Handler(w http.ResponseWriter, r *http.Request) {
	code := http.StatusNotFound
	respondWithError(w, code, http.StatusText(code))
}

func (a *App) HourlyForecast(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	var checkHours string

	hasHours, ok := params["hours"]
	if !ok {
		checkHours = "12"
	} else {
		checkHours = hasHours[0]
	}

	city, state, hours, err := SanitizeHourlyInputs(
		params["city"][0],
		params["state"][0],
		checkHours,
	)
	if err != nil {
		log.Printf("Error checking client inputs: %s\n", err.Error())
		code := http.StatusBadRequest
		respondWithError(w, code, http.StatusText(code))
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
					log.Printf("Error scanning lat/lng from the database: %s\n", err.Error())
					code := http.StatusInternalServerError
					respondWithError(w, code, http.StatusText(code))
				}
			} else {
				a.IncrementLocation(l)
			}
			forecastURL := FetchHourlyForecastURL(l)
			log.Printf("%s", forecastURL)
			forecasts := FetchForecasts(forecastURL)
			hourlyForecasts = a.CacheHourlyForecasts(city, state, hours, forecasts)
		} else if err != nil {
			log.Printf("Error looking up hourly forecasts: %s\n", err.Error())
			code := http.StatusInternalServerError
			respondWithError(w, code, http.StatusText(code))
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
	log.Printf("Scheme is: %s\n", r.URL.Scheme)
	log.Printf("Host is: %s\n", r.URL.Host)
	log.Printf("Path is: %s\n", r.URL.Path)
	log.Printf("Query is: %s\n", r.URL.RawQuery)

	params := r.URL.Query()

	var checkPeriod string

	hasPeriod, ok := params["period"]
	if !ok {
		checkPeriod = "today"
	} else {
		checkPeriod = hasPeriod[0]
	}

	city, state, period, err := SanitizeDetailedInputs(
		params["city"][0],
		params["state"][0],
		checkPeriod,
	)
	// l := Location{City: city.asName, State: state.asName}
	if err != nil {
		log.Printf("Error sanitizing client inputs: %s\n", err.Error())
		code := http.StatusBadRequest
		respondWithError(w, code, http.StatusText(code))
	} else {
		l := Location{City: city.asName, State: state.asName}
		var forecast string
		forecast, err := a.LookupDetailedForecast(city, state, period)
		log.Printf("Redis Error is: %s\n", err.Error())
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
					l.SetCoords(FetchCoords(city.asURL, state.asURL))
					if err = a.RegisterLocation(l); err != nil {
						a.IncrementLocation(l)
					}
				default:
					code := http.StatusBadRequest
					respondWithError(w, code, http.StatusText(code))
				}
			} else {
				a.IncrementLocation(l)
			}
			forecastURL := FetchDetailedForecastURL(l)
			forecasts := FetchForecasts(forecastURL)
			forecast = a.CacheDetailedForecasts(city, state, period, forecasts)
		} else if err != nil {
			log.Printf("Error looking up detailed forecast: %s\n", err.Error())
			code := http.StatusInternalServerError
			respondWithError(w, code, http.StatusText(code))
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
	a.Router.HandleFunc("/api/forecast/detailed", a.DetailedForecast).Queries("city", "{city:[a-zA-Z+]+}", "state", "{state:[a-zA-Z+]+}", "period", "{period:[a-zA-Z+]+}").Methods("GET")
	a.Router.HandleFunc("/api/forecast/detailed", a.DetailedForecast).Queries("city", "{city:[a-zA-Z+]+}", "state", "{state:[a-zA-Z+]+}").Methods("GET")
	a.Router.HandleFunc("/api/forecast/detailed/random", a.RandomDetailedForecast).Methods("GET")
	a.Router.HandleFunc("/api/forecast/hourly", a.HourlyForecast).Queries("city", "{city:[a-zA-Z+]+}", "state", "{state:[a-zA-Z+]+}", "hours", "{hours:[0-9]+}").Methods("GET")
	a.Router.HandleFunc("/api/forecast/hourly", a.HourlyForecast).Queries("city", "{city:[a-zA-Z+]+}", "state", "{state:[a-zA-Z+]+}").Methods("GET")
	// a.Router.HandleFunc("/api/forecast/detailed?city={city}&state={state}&period={period}", a.DetailedForecast).Methods("GET")
	// a.Router.HandleFunc("/api/forecast/detailed?city={city}&state={state}", a.DetailedForecast).Methods("GET")
	// a.Router.HandleFunc("/api/forecast/detailed?random", a.RandomDetailedForecast).Methods("GET")
	// a.Router.HandleFunc("/api/forecast/hourly?city={city}&state={state}&hours={hours}", a.HourlyForecast).Methods("GET")
	// a.Router.HandleFunc("/api/forecast/hourly?city={city}&state={state}", a.HourlyForecast).Methods("GET")
	a.Router.NotFoundHandler = http.HandlerFunc(a.Custom404Handler)
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
