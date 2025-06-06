package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/caseymrm/menuet"
	"github.com/joho/godotenv"
)

type Payload struct {
	Data struct {
		Timings struct {
			Lastthird string `json:"Lastthird"`
			Fajr      string `json:"Fajr"`
			Sunrise   string `json:"Sunrise"`
			Dhuhr     string `json:"Dhuhr"`
			Asr       string `json:"Asr"`
			Maghrib   string `json:"Maghrib"`
			Isha      string `json:"Isha"`
		} `json:"timings"`
	} `json:"data"`
}

func rkt() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	city := os.Getenv("CITY")
	country := os.Getenv("COUNTRY")
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
		current, prayerName := getPrayer(p)
		menuet.App().SetMenuState(&menuet.MenuState{
			Title: "🕌 " + prayerName + " " + current,
		})
		time.Sleep(time.Second * 15)
	}
}

func getPrayer(p Payload) (string, string) {
	v := reflect.ValueOf(p.Data.Timings)

	// Find which is next
	latestTime := fmt.Sprint(p.Data.Timings.Lastthird)
	prayerName := "Last Third:"
	for i := 0; i < v.NumField(); i++ {
		if timeCmpr(fmt.Sprint(v.Field(i).Interface()), time.Now().Format("15:04")) {
			latestTime = fmt.Sprint(v.Field(i).Interface())
			prayerName = v.Type().Field(i).Name + ":"
			break
		} else {
			continue
		}
	}

	return latestTime, prayerName
}

func timeCmpr(time, cur string) bool {
	if time[0:2] > cur[0:2] {
		return true
	} else if time[0:2] == cur[0:2] && time[3:] > cur[3:] {
		return true
	} else {
		return false
	}
}

func main() {
	menuet.App().Label = "com.github.hasan-al-khazraji.rkt"
	go rkt()
	go func() {
		for range time.Tick(1 * time.Hour) {
			go rkt()
		}
	}()
	menuet.App().RunApplication()
}
