package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const redisIndex = 0

var env = os.Getenv("APP_ENV")

func main() {
	flag.Parse()

	router := mux.NewRouter()
	RedisInit(redisIndex)
	InfluxInit()
	WebsocketInit(router)

	//Redis
	router.HandleFunc("/sensors", GetSensors).Methods("GET")
	router.HandleFunc("/sensors/{id}", GetSensor).Methods("GET")
	router.HandleFunc("/sensors", CreateSensor).Methods("POST")
	router.HandleFunc("/sensors/{id}", DeleteSensor).Methods("DELETE")

	//Influx
	router.HandleFunc("/probes", InsertProbe).Methods("POST")
	router.HandleFunc("/probes/query", QueryProbes).Methods("POST")
	router.HandleFunc("/telemetry", InsertTelemetry).Methods("POST")
	router.HandleFunc("/telemetry/query", QueryTelemetry).Methods("POST")

	headers := handlers.AllowedHeaders([]string{"Content-Type", "X-Requested-With"})
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE"})

	srv := &http.Server{
		Handler:      handlers.CORS(headers, origins, methods)(router),
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	if env == "production" {
		log.Println("Running api server in production mode. Port: 8000")
	} else {
		log.Println("Running api server in dev mode. Port: 8000")
	}

	log.Fatal(srv.ListenAndServe())
}
