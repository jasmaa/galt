package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var port = 8000

func main() {

	/*
		// Connect to db
		pgUser := os.Getenv("POSTGRES_USER")
		pgPassword := os.Getenv("POSTGRES_PASSWORD")
		pgHost := os.Getenv("POSTGRES_HOST")
		pgPort := os.Getenv("POSTGRES_PORT")
		pgDB := os.Getenv("POSTGRES_DB")

		connStr := fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			pgUser, pgPassword, pgHost, pgPort, pgDB,
		)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal(err)
		}
	*/

	// Temp logging setup
	l, err := zap.NewDevelopment()
	if err != nil {
		panic("Could not create logger")
	}
	zap.ReplaceGlobals(l)

	// Routes
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hi there")
	}).Methods("GET")

	zap.S().Infof("Server started on port: %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
