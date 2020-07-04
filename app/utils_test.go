package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestSanitizeState(t *testing.T) {
	checkState, _ := sanitizeState("South Dakota")

	target := State{asURL: "SD", asKey: "sd", asName: "SD"}

	if checkState != target {
		t.Errorf("State was incorrect, got: %v, want: %v", checkState, target)
	}
}

func TestSanitizeNotAState(t *testing.T) {
	checkState, err := sanitizeState("West Dakota")

	if err.Error() != "Invalid state name." {
		t.Errorf("Error was incorrect, got: %v, want: Invalid state name.", err)
	}

	target := State{}

	if checkState != target {
		t.Errorf("State was incorrect, got :%v, want: %v", checkState, target)
	}
}

func TestSanitizeCity(t *testing.T) {
	checkCity := sanitizeCity("Salt+Lake CITY")

	target := City{asURL: "Salt+Lake+CITY", asKey: "salt_lake_city", asName: "Salt Lake CITY"}

	if checkCity != target {
		t.Errorf("City was incorrect, got: %v, wanted %v", checkCity, target)
	}
}

func TestSanitizePeriodRelativeDate(t *testing.T) {
	checkPeriod, _ := sanitizePeriod("Tomorrow Night")

	relDate := fmt.Sprintf("%s", strings.ToLower(time.Now().UTC().AddDate(0, 0, 1).Weekday().String()))

	target := Period{
		asKey:     fmt.Sprintf("%s_night", strings.ToLower(relDate)),
		asName:    fmt.Sprintf("%s night", strings.Title(relDate)),
		dayOfWeek: strings.Title(relDate),
		isDaytime: false}

	if checkPeriod != target {
		t.Errorf("Period was incorrect, got: %v, wanted: %v", checkPeriod, target)
	}
}

func TestSanitizePeriodAbsoluteDate(t *testing.T) {
	checkPeriod, _ := sanitizePeriod("WEDNESDAY")

	target := Period{
		asKey:     "wednesday",
		asName:    "Wednesday",
		dayOfWeek: "Wednesday",
		isDaytime: true}

	if checkPeriod != target {
		t.Errorf("Period was incorrect, got: %v, wanted: %v", checkPeriod, target)
	}
}

func TestSanitizePeriodInvalidPeriod(t *testing.T) {
	checkPeriod, err := sanitizePeriod("Reindeer")

	target := Period{}

	if err.Error() != "Invalid period." {
		t.Errorf("Error was incorrect, got: %s, wanted: 'Invalid period.'", err.Error())
	}

	if checkPeriod != target {
		t.Errorf("Period was invalid, got: %v, wanted: %v", checkPeriod, target)
	}
}
