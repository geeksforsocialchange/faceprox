package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	faceloader "github.com/geeksforsocialchange/faceloader/parser"
	"github.com/gorilla/mux"
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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/events/{key}.json", EventDataHandler)
	r.HandleFunc("/events/{key}", EventHandler)
	r.HandleFunc("/page/{key}.json", PageDataHandler)
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
	var result map[string]interface{}

	vars := mux.Vars(r)
	eventUrl := fmt.Sprintf("https://mbasic.facebook.com/events/%v", vars["key"])
	result, err := faceloader.InterfaceFromMbasic(eventUrl)
	if err != nil {
		log.Println(err)
	}
	json.NewEncoder(w).Encode(result)
}

func PageDataHandler(w http.ResponseWriter, r *http.Request) {
	var events []interface{}

	vars := mux.Vars(r)
	pageName := vars["key"]

	mbasicLinks, err := faceloader.GetFacebookEventLinks(pageName)
	mbasicLinks = faceloader.RemoveDuplicateStr(mbasicLinks)

	if err != nil {
		log.Println(err)
	}
	for i := range mbasicLinks {
		log.Println(i, mbasicLinks[i])
		e, err := faceloader.InterfaceFromMbasic(mbasicLinks[i])
		if err != nil {
			log.Println(err)
		}
		events = append(events, e)
	}
	json.NewEncoder(w).Encode(events)
}
