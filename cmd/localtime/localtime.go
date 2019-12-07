package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type response struct {
	Timezone string `json:"timezone"`
}

const (
	url = "http://ip-api.com/json/"
)

var (
	format = flag.String(
		"format", "15:04", "time layout format, in Go's time format style",
	)
	tz = flag.String("tz", "", "timezone in tz data format")
)

func main() {
	flag.Parse()

	var err error
	if *tz == "" {
		tz, err = getTZ()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
	}

	loc, err := time.LoadLocation(*tz)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load tz data: %v\n", err)
		return
	}

	fmt.Println(time.Now().In(loc).Format(*format))
}

func getTZ() (*string, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get ip geolocation: %v\n", err)
	}

	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v\n", err)
	}

	var r response
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %v\n", err)
	}

	return &r.Timezone, nil
}
