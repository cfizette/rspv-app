package main

import (
	"fmt"
	"net/http"
	"rspv-app/controllers"
	"rspv-app/models"
	"rspv-app/views"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Fibonacci112358"
	dbname   = "rsvp"
)

func main() {
	// eventC := controllers.NewEvents()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	es, err := models.NewEventService(psqlInfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer es.Close()
	es.AutoMigrate()
	// es.DestructiveReset()

	home := views.NewView("bootstrap", "static/home")
	create := controllers.NewEvents(es)

	r := mux.NewRouter()

	r.Handle("/", home).Methods("GET")
	r.HandleFunc("/create", create.New).Methods("GET")
	r.HandleFunc("/create", create.Create).Methods("POST")
	r.HandleFunc("/{displayID}", create.ViewEvent).Methods("GET")
	r.HandleFunc("/{displayID}/guest", create.AddGuest).Methods("POST")

	http.ListenAndServe(":3000", r)

	/*
		r := mux.NewRouter()
		r.Handle("/", staticC.Home).Methods("GET")
		r.HandleFunc("/create", eventC.New).Methods("GET")
		r.HandleFunc("/create", eventC.Create).Methods("POST")
		r.HandleFunc("/events", eventC.New).Methods("GET")
		r.HandleFunc("/events", eventC.Attend).Methods("PATCH")

		http.ListenAndServe(":3000", r)
	*/

	/*
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)

		es, err := models.NewEventService(psqlInfo)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer es.Close()
		es.AutoMigrate()
		es.DestructiveReset()

		event := models.Event{
			Name: "Party",
		}

		a := models.Attendee{
			Name: "Billy",
		}

		err = es.Create(&event)
		if err != nil {
			fmt.Println(err)
		}

		err = es.AddAttendee(&event, &a)
		if err != nil {
			fmt.Println(err)
		}

	*/
}
