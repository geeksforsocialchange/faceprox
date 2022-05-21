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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/events/{key}.json", EventDataHandler)
	r.HandleFunc("/events/{key}", EventHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, indexPage)
}

func EventHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, eventPage)
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
