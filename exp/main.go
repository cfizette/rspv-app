package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Fibonacci112358"
	dbname   = "rsvp"
)

type Event struct {
	gorm.Model
	Name      string
	Attendees []Attendee
}

type Attendee struct {
	gorm.Model
	Name    string
	EventID uint
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Error connecting to database")
		fmt.Println(err)
		return
	}
	defer db.Close()

	db.AutoMigrate(&Event{})
	db.AutoMigrate(&Attendee{})

}
