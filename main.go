package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/caseymrm/menuet"
	"github.com/joho/godotenv"
)

type Payload struct {
	Data struct {
		Timings struct {
			Fajr      string `json:"Fajr"`
			Sunrise   string `json:"Sunrise"`
			Dhuhr     string `json:"Dhuhr"`
			Asr       string `json:"Asr"`
			Maghrib   string `json:"Maghrib"`
			Isha      string `json:"Isha"`
			Lastthird string `json:"Lastthird"`
		} `json:"timings"`
	} `json:"data"`
}

func rkt(city, country string) {
	res, err := http.Get("https://api.aladhan.com/v1/timingsByCity/" + time.Now().Format("2-01-2006") + "?city=" + city + "&country=" + country)
	if err != nil {
		return
	}
	defer res.Body.Close()

	var p Payload
	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		log.Fatal(err)
	}
	for {
		menuet.App().SetMenuState(&menuet.MenuState{
			Title: "🕌 " + "Fajr: " + p.Data.Timings.Fajr,
		})
		time.Sleep(time.Second)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	city := os.Getenv("CITY")
	country := os.Getenv("COUNTRY")

	go rkt(city, country)

	menuet.App().RunApplication()
}
