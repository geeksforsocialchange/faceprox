package main

import (
	"alicekaerast/faceprox/lib"
	_ "embed"
	"encoding/json"
	"fmt"
	ics "github.com/arran4/golang-ical"
	faceloader "github.com/geeksforsocialchange/faceloader/parser"
	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
	"log"
	"net/http"
	"time"
)

//go:embed static/event.html
var eventPage string

//go:embed static/index.html
var indexPage string

//go:embed static/events.html
var eventsPage string

var c *cache.Cache

var m *faceloader.MBasic

func main() {
	c = cache.New(5*time.Minute, 10*time.Minute)

	m = faceloader.NewMBasicClient()

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/events/{key}.json", EventDataHandler)
	r.HandleFunc("/events/{key}.ics", EventIcalHandler)
	r.HandleFunc("/events/{key}", EventHandler)
	r.HandleFunc("/page/{key}.json", PageDataHandler)
	r.HandleFunc("/page/{key}.ics", PageIcalHandler)
	r.HandleFunc("/page/{key}", PageHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Listening on http://0.0.0.0:8000")
	log.Fatal(srv.ListenAndServe())
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, indexPage)
}

func EventHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, eventPage)
}

func PageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, eventsPage)
}

func EventDataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	url := fmt.Sprintf("https://mbasic.facebook.com/events/%v", vars["key"])
	result := lib.GetEvent(url, c, m)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func EventIcalHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	url := fmt.Sprintf("https://mbasic.facebook.com/events/%v", vars["key"])
	result := lib.GetEvent(url, c, m)
	ics, err := faceloader.InterfaceToIcal(result)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "calendar/text")
	fmt.Fprint(w, ics.Serialize())
}

func PageIcalHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	events := lib.GetEvents(vars["key"], c, m)
	cal := ics.NewCalendar()

	for i := range events {
		event, err := faceloader.InterfaceToIcal(events[i])
		if err != nil {
			log.Println(err)
		}
		cal.Components = append(cal.Components, &event)
	}
	w.Header().Set("Content-Type", "calendar/text")
	fmt.Fprint(w, cal.Serialize())
}

func PageDataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	events := lib.GetEvents(vars["key"], c, m)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
