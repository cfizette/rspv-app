package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

var (
	ErrInvalidID = errors.New("models: Invalid ID provided")
)

type Event struct {
	gorm.Model
	Name      string
	Guests    []Guest `gorm:"foreignkey:EventDisplayID;association_foreignkey:DisplayID"`
	DisplayID *string `gorm:"not null;unique_index"`
}

type EventService struct {
	db *gorm.DB
}

func NewEventService(connectionInfo string) (*EventService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &EventService{
		db: db,
	}, nil
}

// Close closes the database connection
func (es *EventService) Close() error {
	return es.db.Close()
}

// Create creates a new event
func (es *EventService) Create(e *Event) error {
	return es.db.Create(e).Error
}

// Update updates the provided event with the provided data
func (es *EventService) Update(e *Event) error {
	return es.db.Save(e).Error
}

// Delete will delete the event with the provided ID
func (es *EventService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	e := Event{Model: gorm.Model{ID: id}}
	return es.db.Delete(&e).Error
}

// ByDisplayID will return the event with the provided displayID
func (es *EventService) ByDisplayID(id string) (*Event, error) {
	var e Event
	db := es.db.Where("display_id = ?", id)
	err := db.First(&e).Error
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// AddGuest creates an association between the given Event and Attendee
func (es *EventService) AddGuest(e *Event, g *Guest) error {
	// Currently possible to add gest to deleted event
	err := es.db.Model(e).Association("Guests").Append(*g).Error
	return err
}

// GetGuests returns an array of all Guests associates with an Event.
//TODO: look into using pointers to return Guests
func (es *EventService) GetGuests(e *Event) ([]Guest, error) {
	var guests []Guest
	err := es.db.Where("event_display_id = ?", e.DisplayID).Find(&guests).Error
	return guests, err
}

// DestructiveReset drops the events and attendees table and rebuilds it
func (es *EventService) DestructiveReset() error {
	err := es.db.DropTableIfExists(&Event{}).Error
	if err != nil {
		return err
	}
	err = es.db.DropTableIfExists(&Guest{}).Error
	if err != nil {
		return err
	}
	return es.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate the
// event and attendee table
func (es *EventService) AutoMigrate() error {
	if err := es.db.AutoMigrate(&Event{}).Error; err != nil {
		return err
	}
	if err := es.db.AutoMigrate(&Guest{}).Error; err != nil {
		return err
	}
	return nil
}
