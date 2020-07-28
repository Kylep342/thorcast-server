package app

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

	_ "github.com/jackc/pgx/stdlib"
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

// App contains necessary components to run the webserver
// Router is a pointer to a mux Router
// Logger is an http handler
// DB is a pointer to a db
// Redis is a pointer to a redis client
type App struct {
	Router *mux.Router
	Logger http.Handler
	DB     *sql.DB
	Redis  *redis.Client
}

// InitializeRoutes creates all endpoints for the api
func (a *App) InitializeRoutes() {
	a.Router.HandleFunc("/api/forecast/detailed", a.DetailedForecastHandler).Queries("city", "{city:[a-zA-Z+]+}", "state", "{state:[a-zA-Z+]+}", "period", "{period:[a-zA-Z+]+}").Methods("GET")
	a.Router.HandleFunc("/api/forecast/detailed", a.DetailedForecastHandler).Queries("city", "{city:[a-zA-Z+]+}", "state", "{state:[a-zA-Z+]+}").Methods("GET")
	a.Router.HandleFunc("/api/forecast/detailed/random", a.RandomDetailedForecastHandler).Methods("GET")
	a.Router.HandleFunc("/api/forecast/hourly", a.HourlyForecastHandler).Queries("city", "{city:[a-zA-Z+]+}", "state", "{state:[a-zA-Z+]+}", "hours", "{hours:[0-9]+}").Methods("GET")
	a.Router.HandleFunc("/api/forecast/hourly", a.HourlyForecastHandler).Queries("city", "{city:[a-zA-Z+]+}", "state", "{state:[a-zA-Z+]+}").Methods("GET")
	a.Router.NotFoundHandler = http.HandlerFunc(a.Custom404Handler)
}

// Initialize creates the application as a whole
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
	a.DB, err = sql.Open("pgx", sqlDataSource)
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

// Run starts the app to listen on the port specitied by the env variable SERVER_PORT
func (a *App) Run() {
	port := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	log.Fatal(http.ListenAndServe(port, a.Logger))
}
