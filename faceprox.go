package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
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
	var result []map[string]interface{}

	vars := mux.Vars(r)
	eventUrl := fmt.Sprintf("https://mbasic.facebook.com/events/%v", vars["key"])
	res, _ := http.Get(eventUrl)
	doc, _ := goquery.NewDocumentFromReader(res.Body)

	selector := `script[type="application/ld+json"]`
	scripts := doc.Find(selector)
	scripts.Each(func(i int, s *goquery.Selection) {
		var decoded map[string]interface{}
		text := s.Text()
		text = strings.Replace(text, "//<![CDATA[", "", -1)
		text = strings.Replace(text, "//]]", "", -1)
		text = strings.Replace(text, ">", "", -1)
		err := json.Unmarshal([]byte(text), &decoded)
		if err != nil {
			log.Println(err)
		}
		if err == nil {
			result = append(result, decoded)
		}
	})

	json.NewEncoder(w).Encode(result)
}
