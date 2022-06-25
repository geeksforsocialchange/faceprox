package lib

import (
	faceloader "github.com/geeksforsocialchange/faceloader/parser"
	"github.com/patrickmn/go-cache"
	"log"
)

func GetEvent(e string, c *cache.Cache) map[string]interface{} {
	var result map[string]interface{}
	var err error

	cachedResult, found := c.Get(e)
	if found {
		log.Println("cache hit")
		result = cachedResult.(map[string]interface{})
	} else {
		log.Println("cache miss")
		result, err = faceloader.InterfaceFromMbasic(e)
		if err != nil {
			log.Println(err)
		}
		c.Set(e, result, cache.DefaultExpiration)
	}
	return result
}

func GetLinks(p string, c *cache.Cache) []string {
	var mbasicLinks []string
	var err error

	cachedLinks, found := c.Get(p)
	if found {
		log.Println("cache hit")
		mbasicLinks = cachedLinks.([]string)
	} else {
		log.Println("cache miss")
		mbasicLinks, err = faceloader.GetFacebookEventLinks(p)
		mbasicLinks = faceloader.RemoveDuplicateStr(mbasicLinks)
		if err != nil {
			log.Println(err)
		} else {
			c.Set(p, mbasicLinks, cache.DefaultExpiration)
		}
	}
	return mbasicLinks
}

func GetEvents(p string, c *cache.Cache) []map[string]interface{} {
	var events []map[string]interface{}

	l := GetLinks(p, c)

	for i := range l {
		event := GetEvent(l[i], c)
		events = append(events, event)
	}
	return events
}
