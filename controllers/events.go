package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"rspv-app/models"
	"rspv-app/views"

	"github.com/gorilla/mux"

	"github.com/gorilla/schema"
)

type Events struct {
	NewView     *views.View
	CreatedView *views.View
	EventView   *views.View
	es          *models.EventService
}

type EventDetails struct {
	Event  *models.Event
	Guests []models.Guest
}

// NewEvents is used to create a new Events controller
// This will panic if the templates are not parsed correctly
// and should only be used during setup.
func NewEvents(es *models.EventService) *Events {
	return &Events{
		NewView:     views.NewView("bootstrap", "events/new"),
		CreatedView: views.NewView("bootstrap", "events/created"),
		EventView:   views.NewView("bootstrap", "events/view"),
		es:          es,
	}
}

// EventCreationForm contains all the information required
// to create a new Event
type EventCreationForm struct {
	Name string `schema:"name"`
}

// AddGuestForm contains the information of the guest to be
// added as well as the event to be added to
type AddGuestForm struct {
	Name string `schema:"name"`
}

// GET /create
func (e *Events) New(w http.ResponseWriter, r *http.Request) {
	if err := e.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// POST /create
func (e *Events) Create(w http.ResponseWriter, r *http.Request) {
	var form EventCreationForm
	if err := parseForm(&form, r); err != nil {
		panic(err)
	}
	dispID := randomString(10)
	event := models.Event{
		Name:      form.Name,
		DisplayID: &dispID,
	}
	if err := e.es.Create(&event); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := e.CreatedView.Render(w, event); err != nil {
		panic(err)
	}
}

// ViewEvent renders the view for an Event and its Guests
func (e *Events) ViewEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	displayID := vars["displayID"]
	event, err := e.es.ByDisplayID(displayID)
	if err != nil {
		panic(err)
	}
	guests, err := e.es.GetGuests(event)
	if err != nil {
		panic(err)
	}
	eventDetails := EventDetails{
		Event:  event,
		Guests: guests,
	}
	if err := e.EventView.Render(w, eventDetails); err != nil {
		panic(err)
	}

}

func (e *Events) AddGuest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	displayID := vars["displayID"]
	var form AddGuestForm
	if err := parseForm(&form, r); err != nil {
		panic(err)
	}
	event := models.Event{
		DisplayID: &displayID,
	}
	guest := models.Guest{
		Name: &form.Name,
	}
	fmt.Println(*event.DisplayID)
	if err := e.es.AddGuest(&event, &guest); err != nil {
		panic(err)
	}
}

func parseForm(dst interface{}, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	dec := schema.NewDecoder()
	if err := dec.Decode(dst, r.PostForm); err != nil {
		return err
	}
	return nil
}

// TODO: Fix this since not truely random
// When the server is restarted, ids are created in
// the same order as before....
func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(97, 122))
	}
	return string(bytes)
}
