package main

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

var timesOfDay = []string{"", "_night"}
var separatorRE = regexp.MustCompile(`[_+ ]+`)
var periodRE = regexp.MustCompile(`(?i)(sunday|monday|tuesday|wednesday|thursday|friday|saturday)( night)?`)
var relDateRE = regexp.MustCompile(`(?i)(today|tonight|tomorrow)( night)?`)
var stateCodes = map[string]string{
	"alabama": "AL", "al": "AL",
	"alaska": "AK", "ak": "AK",
	"arizona": "AZ", "az": "AZ",
	"arkansas": "AR", "ar": "AR",
	"california": "CA", "ca": "CA",
	"colorado": "CO", "co": "CO",
	"connecticut": "CT", "ct": "CT",
	"delaware": "DE", "de": "DE",
	"florida": "FL", "fl": "FL",
	"georgia": "GA", "ga": "GA",
	"hawaii": "HI", "hi": "HI",
	"idaho": "ID", "id": "ID",
	"illinois": "IL", "il": "IL",
	"indiana": "IN", "in": "IN",
	"iowa": "IA", "ia": "IA",
	"kansas": "KS", "ks": "KS",
	"kentucky": "KY", "ky": "KY",
	"louisiana": "LA", "la": "LA",
	"maine": "ME", "me": "ME",
	"maryland": "MD", "md": "MD",
	"massachusetts": "MA", "ma": "MA",
	"michigan": "MI", "mi": "MI",
	"minnesota": "MN", "mn": "MN",
	"mississippi": "MS", "ms": "MS",
	"missouri": "MO", "mo": "MO",
	"montana": "MT", "mt": "MT",
	"nebraska": "NE", "ne": "NE",
	"nevada": "NV", "nv": "NV",
	"new hampshire": "NH", "nh": "NH",
	"new jersey": "NJ", "nj": "NJ",
	"new mexico": "NM", "nm": "NM",
	"new york": "NY", "ny": "NY",
	"north carolina": "NC", "nc": "NC",
	"north dakota": "ND", "nd": "ND",
	"ohio": "OH", "oh": "OH",
	"oklahoma": "OK", "ok": "OK",
	"oregon": "OR", "or": "OR",
	"pennsylvania": "PA", "pa": "PA",
	"rhode island": "RI", "ri": "RI",
	"south carolina": "SC", "sc": "SC",
	"south dakota": "SD", "sd": "SD",
	"tennessee": "TN", "tn": "TN",
	"texas": "TX", "tx": "TX",
	"utah": "UT", "ut": "UT",
	"vermont": "VT", "vt": "VT",
	"virginia": "VA", "va": "VA",
	"washington": "WA", "wa": "WA",
	"west virginia": "WV", "wv": "WV",
	"wisconsin": "WI", "wi": "WI",
	"wyoming": "WY", "wy": "WY"}

type City struct {
	asURL  string
	asKey  string
	asName string
}

type State struct {
	asURL  string
	asKey  string
	asName string
}

type Period struct {
	asKey     string
	asName 	  string
	dayOfWeek string
	isDaytime bool
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func SanitizeInputs(city string, state string, period string) (City, State, Period, error) {
	cleanState, err := sanitizeState(state)
	if err != nil {
		return City{}, State{}, Period{}, err
	}
	cleanPeriod, err := sanitizePeriod(period)
	if err != nil {
		return City{}, State{}, Period{}, err
	}
	cleanCity := sanitizeCity(city)
	return cleanCity, cleanState, cleanPeriod, nil
}

func sanitizeCity(city string) City {
	return City{asURL: separatorRE.ReplaceAllString(city, "+"),
		asKey:  strings.ToLower(separatorRE.ReplaceAllString(city, "_")),
		asName: separatorRE.ReplaceAllString(city, " ")}
}

func sanitizeState(state string) (State, error) {
	key := strings.ToLower(separatorRE.ReplaceAllString(state, " "))
	if cleanState, ok := stateCodes[key]; ok {
		return State{asURL: cleanState, asKey: strings.ToLower(cleanState), asName: cleanState}, nil
	}
	return State{}, errors.New("Invalid state name.")
}

func sanitizePeriod(period string) (Period, error) {
	var cleanPeriod string
	switch {
	case strings.Contains(strings.ToLower(period), "today"):
		cleanPeriod = time.Now().Weekday().String()
	case strings.Contains(strings.ToLower(period), "tonight"):
		cleanPeriod = fmt.Sprintf("%s night", time.Now().Weekday().String())
	case strings.Contains(strings.ToLower(period), "tomorrow"):
		cleanPeriod = time.Now().AddDate(0,0,1).Weekday().String()
	case strings.Contains(strings.ToLower(period), "tomorrow night"):
		cleanPeriod = fmt.Sprintf("%s night", time.Now().AddDate(0,0,1).Weekday().String())
	default:
		cleanPeriod = period
	}
	m := periodRE.FindStringSubmatch(separatorRE.ReplaceAllString(cleanPeriod, " "))
	if m != nil && m[2] == "" {
		return Period{
			asKey: strings.ToLower(m[1]),
			asName: strings.Title(m[1]),
			dayOfWeek: strings.Title(m[1]),
			isDaytime: true}, nil
	} else if m != nil {
		return Period{
			asKey: fmt.Sprintf(
				"%s_%s",
				strings.ToLower(m[1]),
				strings.ToLower(m[2])),
			asName: fmt.Sprintf(
				"%s %s",
				strings.Title(m[1]),
				strings.ToLower(m[2])),
			dayOfWeek: strings.Title(m[1]),
			isDaytime: false}, nil
	} else {
		return Period{}, errors.New("Invalid period.")
	}
}

func randomPeriod() Period {
	dayOfWeek := time.Now().AddDate(0, 0, rand.Intn(7)).Weekday().String()
	timeOfDay := timesOfDay[rand.Intn(2)]
	p, _ := sanitizePeriod(fmt.Sprintf("%s%s", dayOfWeek, timeOfDay))
	return p
}