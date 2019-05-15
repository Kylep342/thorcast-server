package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type config struct {
	username string
	password string
	host     string
	port     string
	dbName   string
}

func (conf *config) configure() {
	conf.username = os.Getenv("THORCAST_DB_USERNAME")
	conf.password = os.Getenv("THORCAST_DB_PASSWORD")
	conf.host = os.Getenv("THORCAST_DB_HOST")
	conf.port = os.Getenv("THORCAST_DB_PORT")
	conf.dbName = os.Getenv("THORCAST_DB_NAME")
}

var dbConf = config{}

type App struct {
	Router *mux.Router
	Logger http.Handler
	DB     *sql.DB
}

func (a *App) LookupFC(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	city := strings.ReplaceAll(vars["city"], "+", " ")
	state := strings.ReplaceAll(vars["state"], "+", " ")
	period := strings.ReplaceAll(vars["period"], "+", " ")

	lCity, lState, lPeriod, err := SanitizeInputs(city, state, period)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		l := Location{City: lCity.asName, State: lState.asName}
		row := a.DB.QueryRow("SELECT lat, lng FROM geocodex WHERE LOWER(city) = LOWER($1) AND state = $2;", l.City, l.State)
		if err := row.Scan(&l.Lat, &l.Lng); err != nil {
			switch err {
			case sql.ErrNoRows:
				l.SetCoords(FetchCoords(lCity.asURL, lState.asURL))
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
		forecast := SelectForecast(forecasts, lPeriod)
		resp := map[string]string{"detailedForecast": forecast}
		respondWithJSON(w, http.StatusOK, resp)
	}
}

func (a *App) InitializeRoutes() {
	a.Router.HandleFunc("/api/forecast/city={city}&state={state}&period={period}", a.LookupFC).Methods("GET")
}

func (a *App) Initialize() {
	var err error
	dbConf.configure()
	dataSource := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbConf.username,
		dbConf.password,
		dbConf.host,
		dbConf.port,
		dbConf.dbName)
	a.DB, err = sql.Open("postgres", dataSource)
	if err != nil {
		log.Fatal(err)
	}
	a.Router = mux.NewRouter()
	a.Logger = handlers.CombinedLoggingHandler(os.Stdout, a.Router)
	a.InitializeRoutes()
}

func (a *App) Run() {
	log.Fatal(http.ListenAndServe(":8000", a.Logger))
}
