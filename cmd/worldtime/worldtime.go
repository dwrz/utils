package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Location struct {
	Name     string
	Timezone string
}

var locations = []Location{
	{Name: "NYC     ", Timezone: "America/New_York"},
	{Name: "Napoli  ", Timezone: "Europe/Rome"},
	{Name: "Shanghai", Timezone: "Asia/Shanghai"},
	{Name: "UTC     ", Timezone: "UTC"},
}

const (
	dateFormat = "20060102"
	timeFormat = "150405 -0700 MST"
)

func main() {
	now := time.Now()

	var str strings.Builder
	for _, loc := range locations {
		tz, err := time.LoadLocation(loc.Timezone)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}

		lt := now.In(tz)

		str.WriteString(fmt.Sprintf(
			"%s %s %d/365 %d/52 %d/7 %s\n",
			loc.Name,
			lt.Format(dateFormat),
			lt.YearDay(),
			(lt.YearDay()/7)+1,
			int(lt.Weekday()),
			lt.Format(timeFormat),
		))
	}

	fmt.Println(str.String())
}
