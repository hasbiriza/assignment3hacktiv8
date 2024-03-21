package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Weather struct {
	Water int `json:"water"` // dalam satuan meter
	Wind  int `json:"wind"`  // dalam satuan meter per detik
}

func (w *Weather) checkStatus() (resWater string, resWind string) {
	switch {
	case w.Water < 5:
		resWater = "Aman"
	case w.Water >= 6 && w.Water <= 8:
		resWater = "Siaga"
	default:
		resWater = "Bahaya"
	}

	switch {
	case w.Wind < 6:
		resWind = "Aman"
	case w.Wind >= 7 && w.Wind <= 15:
		resWind = "Siaga"
	default:
		resWind = "Bahaya"
	}
	return resWater, resWind
}

func generateJSON() {
	weather := Weather{
		Water: rand.Intn(100) + 1,
		Wind:  rand.Intn(100) + 1,
	}

	file, err := os.Create("weather.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(weather)
	if err != nil {
		log.Fatal(err)
	}
}

func updateJSONEvery(duration time.Duration) {
	for range time.Tick(duration) {
		generateJSON()
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("weather.json")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer file.Close()

	var weather Weather
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&weather)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	resWater, resWind := weather.checkStatus()
	status := fmt.Sprintf("Status Air: %s, Status Angin: %s", resWater, resWind)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func main() {
	http.HandleFunc("/status", statusHandler)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	go updateJSONEvery(15 * time.Second)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
