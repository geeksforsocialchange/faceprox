package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
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

func main() {
	c = cache.New(5*time.Minute, 10*time.Minute)

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
	var err error

	vars := mux.Vars(r)
	eventUrl := fmt.Sprintf("https://mbasic.facebook.com/events/%v", vars["key"])
	cachedResult, found := c.Get(eventUrl)
	if found {
		log.Println("cache hit")
		result = cachedResult.(map[string]interface{})
	} else {
		log.Println("cache miss")
		result, err = faceloader.InterfaceFromMbasic(eventUrl)
		if err != nil {
			log.Println(err)
		}
		c.Set(eventUrl, result, cache.DefaultExpiration)
	}
	json.NewEncoder(w).Encode(result)
}

func PageDataHandler(w http.ResponseWriter, r *http.Request) {
	var events []interface{}
	var mbasicLinks []string
	var err error

	vars := mux.Vars(r)
	pageName := vars["key"]

	cachedLinks, found := c.Get(pageName)
	if found {
		log.Println("cache hit")
		mbasicLinks = cachedLinks.([]string)
	} else {
		log.Println("cache miss")
		mbasicLinks, err = faceloader.GetFacebookEventLinks(pageName)
		mbasicLinks = faceloader.RemoveDuplicateStr(mbasicLinks)
		if err != nil {
			log.Println(err)
		} else {
			c.Set(pageName, mbasicLinks, cache.DefaultExpiration)
		}
	}
	for i := range mbasicLinks {
		link := mbasicLinks[i]

		event, found := c.Get(link)
		if found {
			log.Println("cache hit")
			events = append(events, event.(interface{}))
		} else {
			log.Println("cache miss")
			e, err := faceloader.InterfaceFromMbasic(mbasicLinks[i])
			if err != nil {
				log.Println(err)
			} else {
				c.Set(link, e, cache.DefaultExpiration)
			}
			events = append(events, e)
		}

	}
	json.NewEncoder(w).Encode(events)
}
