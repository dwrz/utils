package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

const (
	SECONDS_IN_DAY  = 86400
	SECONDS_IN_WEEK = 604800
)

var (
	now = time.Now()

	year  = flag.Int("year", 0, "birth year")
	month = flag.Int("month", 0, "birth month")
	day   = flag.Int("day", 0, "birth day")

	date = flag.String(
		"date",
		"",
		"the RFC3339 date and time from which to measure age",
	)
	output = flag.String(
		"output",
		"years",
		"the duration to output [years|months|weeks|days|seconds]",
	)
)

func main() {
	flag.Parse()

	if *year <= 0 {
		fmt.Fprintf(os.Stderr, "invalid year: %d\n", *year)
		return
	}
	if *month <= 0 || *month > 12 {
		fmt.Fprintf(os.Stderr, "invalid month: %d\n", *month)
		return
	}
	if *day <= 0 || *day > 31 {
		fmt.Fprintf(os.Stderr, "invalid day: %d\n", *day)
		return
	}

	if *date != "" {
		var err error
		now, err = time.Parse(time.RFC3339, *date)
		if err != nil {
			fmt.Fprintf(
				os.Stderr, "failed to parse date: %v\n", err,
			)
			return
		}
	}

	birthMonth := time.Date(
		*year, time.Month(*month), 1, 0, 0, 0, 0, time.Local,
	)
	lastDayOfBirthMonth := birthMonth.
		AddDate(0, 1, 0).Add(-time.Nanosecond).Day()
	if *day > lastDayOfBirthMonth {
		fmt.Fprintf(
			os.Stderr,
			"invalid day: %d is after last day of month %d\n",
			*day, lastDayOfBirthMonth,
		)
		return
	}

	birthday := time.Date(
		*year, time.Month(*month), *day, 0, 0, 0, 0, time.Local,
	)

	if now.Before(birthday) {
		fmt.Fprintf(os.Stderr, "cannot measure negative age")
		return
	}

	ageInSeconds := int(now.Sub(birthday).Seconds())
	ageInCalendarYears := now.Year() - birthday.Year()
	birthdayPassed := func() bool {
		switch {
		case ageInCalendarYears > 0 &&
			(now.Month() > birthday.Month()),
			now.Month() == birthday.Month() &&
				now.Day() >= birthday.Day():
			return true
		default:
			return false
		}
	}()

	switch *output {

	case "years":
		if birthdayPassed {
			fmt.Println(ageInCalendarYears)
			return
		}
		fmt.Println(ageInCalendarYears - 1)

	case "months":
		yearsInMonths := (ageInCalendarYears * 12)
		ageInCalendarMonths := int(now.Month() - birthday.Month())
		switch {
		case birthdayPassed && now.Day() >= birthday.Day():
			fmt.Println(yearsInMonths + ageInCalendarMonths)
		case birthdayPassed && now.Day() < birthday.Day():
			fmt.Println(yearsInMonths + ageInCalendarMonths - 1)
		case !birthdayPassed && now.Day() >= birthday.Day():
			if ageInCalendarYears == 0 {
				fmt.Println(ageInCalendarMonths)
				return
			}
			fmt.Println(yearsInMonths - ageInCalendarMonths)
		case !birthdayPassed && now.Day() < birthday.Day():
			if ageInCalendarYears == 0 {
				fmt.Println(ageInCalendarMonths - 1)
				return
			}
			fmt.Println(yearsInMonths - ageInCalendarMonths - 1)
		}

	case "weeks":
		fmt.Println(ageInSeconds / SECONDS_IN_WEEK)

	case "days":
		fmt.Println(ageInSeconds / SECONDS_IN_DAY)

	case "seconds":
		fmt.Println(ageInSeconds)

	default:
		fmt.Fprintf(
			os.Stderr, "unrecognized output flag: %s\n", *output,
		)
		return
	}
}
