package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-redis/redis"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type config struct {
	sqlUsername 	string
	sqlPassword 	string
	sqlHost     	string
	sqlPort     	string
	sqlDbName   	string
	redisPassword	string
	redisHost		string
	redisPort		string
	redisDb			int
}

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

func (a *App) Forecast(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	city, state, period, err := SanitizeInputs(
		vars["city"],
		vars["state"],
		vars["period"],
	)
	l := Location{City: city.asName, State: state.asName}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		var forecast string
		forecast, err := a.LookupForecast(city, state, period)
		if err == redis.Nil {
			row := a.DB.QueryRow(
			"SELECT lat, lng FROM geocodex WHERE LOWER(city) = LOWER($1) AND state = $2;",
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
			forecastURL := FetchForecastURL(l)
			forecasts := FetchForecasts(forecastURL)
			forecast = a.CacheForecasts(city, state, period, forecasts)
		} else if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		} else {
			a.IncrementLocation(l)
		}
		resp := map[string]string{
			"detailedForecast": forecast,
			"city": city.asName,
			"state": state.asName,
			"period": period.asName}
		respondWithJSON(w, http.StatusOK, resp)
	}
}

func (a *App) InitializeRoutes() {
	a.Router.HandleFunc("/api/forecast/city={city}&state={state}&period={period}", a.Forecast).Methods("GET")
}

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
		Addr:		fmt.Sprintf("%s:%s", conf.redisHost, conf.redisPort),
		Password:	conf.redisPassword,
		DB:			conf.redisDb,
	})
	a.Router = mux.NewRouter()
	a.Logger = handlers.CombinedLoggingHandler(os.Stdout, a.Router)
	a.InitializeRoutes()
}

func (a *App) Run() {
	port := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	log.Fatal(http.ListenAndServe(port , a.Logger))
}
