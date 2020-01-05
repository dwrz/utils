package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	now = time.Now()

	year  = flag.Int("year", 0, "year to use; defaults to current year if not set or 0")
	month = flag.Int("month", 0, "month to use; defaults to current month if not set or 0")

	ymMode = flag.Bool("ym", false, "create YYMM directories")
)

func main() {
	flag.Parse()

	if *year < 0 {
		fmt.Fprintf(os.Stderr, "invalid year: %d\n", *year)
		return
	}
	if *year == 0 {
		currentYear := now.Year()
		year = &currentYear
	}

	if *ymMode {
		mkdirYM(*year)
		return
	}

	if *month < 0 || *month > 12 {
		fmt.Fprintf(os.Stderr, "invalid month: %d\n", month)
		return
	}
	if *month == 0 {
		currentMonth := int(now.Month())
		month = &currentMonth
	}

	mkdirYMD(*year, *month)
}

func mkdirYM(year int) {
	for m := 1; m <= 12; m++ {
		name := fmt.Sprintf("%d%02d", year, m)
		if err := os.Mkdir(name, 0777); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}
}

func mkdirYMD(year, month int) {
	lastDayOfMonth := time.Date(
		year, time.Month(month), 1, 0, 0, 0, 0, time.Local,
	).AddDate(0, 1, 0).Add(-time.Nanosecond).Day()

	for d := 1; d <= lastDayOfMonth; d++ {
		name := fmt.Sprintf("%d%02d%02d", year, month, d)
		if err := os.Mkdir(name, 0777); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}
}
