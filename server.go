package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/smahjoub/events-api/handlers"
	"github.com/smahjoub/events-api/store"
)

// Args args used to run the server
type Args struct {
	// postgres connection string, of the form,
	// e.g "postgres://user:password@localhost:5432/database?sslmode=disable"
	conn string
	// port for the server of the form,
	// e.g ":8080"
	port string
}

// Run run the server based on given args
func Run(args Args) error {
	// router
	router := mux.NewRouter().
		PathPrefix("/api/v1/"). // add prefix for v1 api `/api/v1/`
		Subrouter()

	st := store.NewPostgresEventStore(args.conn)
	hnd := handlers.NewEventHandler(st)
	RegisterAllRoutes(router, hnd)

	// start server
	log.Println("Starting server at port: ", args.port)
	return http.ListenAndServe(args.port, router)
}

// RegisterAllRoutes registers all routes of the api
func RegisterAllRoutes(router *mux.Router, hnd handlers.IEventHandler) {

	// set content type
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	// get events
	router.HandleFunc("/event", hnd.Get).Methods(http.MethodGet)
	// create events
	router.HandleFunc("/event", hnd.Create).Methods(http.MethodPost)
	// delete event
	router.HandleFunc("/event", hnd.Delete).Methods(http.MethodDelete)

	// cancel event
	router.HandleFunc("/event/cancel", hnd.Cancel).Methods(http.MethodPatch)
	// update event details
	router.HandleFunc("/event/details", hnd.UpdateDetails).Methods(http.MethodPut)
	// reschedule event
	router.HandleFunc("/event/reschedule", hnd.Reschedule).Methods(http.MethodPatch)

	// list events
	router.HandleFunc("/events", hnd.List).Methods(http.MethodGet)
}
