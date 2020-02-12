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
	Attendees []Attendee
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

// ByID will return the event with the provided ID
func (es *EventService) ByID(id uint) (*Event, error) {
	var e Event
	db := es.db.Where("id = ?", id)
	err := db.First(&e).Error
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// AddAttendee creates an association between the given Event and Attendee
func (es *EventService) AddAttendee(e *Event, a *Attendee) error {
	// Currently possible to add attendee to deleted event
	err := es.db.Model(e).Association("Attendees").Append(*a).Error
	return err
}

// DestructiveReset drops the events and attendees table and rebuilds it
func (es *EventService) DestructiveReset() error {
	err := es.db.DropTableIfExists(&Event{}).Error
	if err != nil {
		return err
	}
	err = es.db.DropTableIfExists(&Attendee{}).Error
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
	if err := es.db.AutoMigrate(&Attendee{}).Error; err != nil {
		return err
	}
	return nil
}
