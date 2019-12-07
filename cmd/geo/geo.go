package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type response struct {
	City        string  `json:"city"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	District    string  `json:"district"`
	Latitude    float32 `json:"lat"`
	Longitude   float32 `json:"lon"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	Timezone    string  `json:"timezone"`
	Zip         string  `json:"zip"`
}

const (
	url = "http://ip-api.com/json/"
)

func main() {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("%v", err)
		return
	}

	defer res.Body.Close()

	var r response
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		fmt.Println("%v\n", err)
		return
	}

	fmt.Printf("%f,%f\n", r.Latitude, r.Longitude)
}
