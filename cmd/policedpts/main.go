package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	forcesListURL    = "https://data.police.uk/api/forces"
	forcesDetailsURL = "https://data.police.uk/api/forces/"
)

type force struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type EngagementMethod struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

type forceDetails struct {
	Telephone  string             `json:"telephone"`
	EngMethods []EngagementMethod `json:"engagement_methods"`
}

func main() {
	resp, err := http.Get(forcesListURL)
	if err != nil {
		log.Fatal(err)
	}
	var forces []force
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()
	err = json.NewDecoder(resp.Body).Decode(&forces)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range forces {
		r, err := http.Get(forcesDetailsURL + f.ID)
		if err != nil {
			log.Println(err)
		}
		details := forceDetails{
			EngMethods: []EngagementMethod{},
		}
		err = json.NewDecoder(r.Body).Decode(&details)
		if err != nil {
			log.Println(err)
		}
		var emFB EngagementMethod
		for _, em := range details.EngMethods {
			if em.Title == "Facebook" {
				emFB = em
				break
			}
		}
		if emFB.Title != "" {
			fmt.Printf("\"%s\",\"%s\",\"%s\"\n", f.Name, details.Telephone, emFB.URL)
		}
	}
}
