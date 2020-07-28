package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Slice containing all string representations of different times of day
// empty string is day time, "_night" is night time
var timesOfDay = []string{"", "[_ ]night"}

// Regex to match any expected separator characters in a string within
// a URL parameter
var separatorRE = regexp.MustCompile(`[_+ ]+`)

// Regex used to match accepted days of the week and times of day
var periodRE = regexp.MustCompile(`(?i)(sunday|monday|tuesday|wednesday|thursday|friday|saturday) ?(night)?`)

// Regex used to match relative days and times (e.g. today, tomorrow night, etc.)
var relDateRE = regexp.MustCompile(`(?i)(today|tonight|tomorrow) ?(night)?`)

// Map containing all accepted state names and abbreviations
// Matches output capitalized 2 character postal codes
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

// City contains a city name in several string representations
// asURL: Case insensitive, separated by plus signs
// asKey: Lowercase, separated by underscores
// asName: Proper case, separated by spaces
type City struct {
	asURL  string
	asKey  string
	asName string
}

// State contains a state name in several string representations
// asURL: Uppercase
// asKey: Lowercase
// asName: Uppercase
type State struct {
	asURL  string
	asKey  string
	asName string
}

// Period contains a time period in several string representations
// as well as a boolean representing the state of daytime in the period
type Period struct {
	asKey     string
	asName    string
	dayOfWeek string
	isDaytime bool
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// SanitizeDetailedInputs is a wrapper function to validate all inputs
// Returns City, State, Period, and nil error on success
func SanitizeDetailedInputs(city string, state string, period string) (City, State, Period, error) {
	cleanState, err := SanitizeState(state)
	if err != nil {
		return City{}, State{}, Period{}, err
	}
	cleanPeriod, err := sanitizePeriod(period)
	if err != nil {
		return City{}, State{}, Period{}, err
	}
	cleanCity := SanitizeCity(city)
	return cleanCity, cleanState, cleanPeriod, nil
}

// SanitizeHourlyInputs is a wrapper function to validate all inputs
// Returns City, State, hours int64, and nil error on success
func SanitizeHourlyInputs(city string, state string, hours string) (City, State, int64, error) {
	cleanCity, cleanState, err := sanitizeLocation(city, state)
	if err != nil {
		return City{}, State{}, 0, err
	}
	cleanHours, err := strconv.ParseInt(hours, 10, 64)
	if err != nil {
		return City{}, State{}, 0, err
	}
	return cleanCity, cleanState, cleanHours, nil
}

// SanitizeCity creates a City struct from a given city name string
func SanitizeCity(city string) City {
	return City{asURL: separatorRE.ReplaceAllString(city, "+"),
		asKey:  strings.ToLower(separatorRE.ReplaceAllString(city, "_")),
		asName: separatorRE.ReplaceAllString(city, " ")}
}

// SanitizeState creates a State struct from a given state name string
func SanitizeState(state string) (State, error) {
	key := strings.ToLower(separatorRE.ReplaceAllString(state, " "))
	if cleanState, ok := stateCodes[key]; ok {
		return State{asURL: cleanState, asKey: strings.ToLower(cleanState), asName: cleanState}, nil
	}
	return State{}, errors.New("invalid state name")
}

func sanitizeLocation(city string, state string) (City, State, error) {
	cleanCity := SanitizeCity(city)
	cleanState, err := SanitizeState(state)
	return cleanCity, cleanState, err
}

// sanitizePeriod creates a Period struct from a given period name string
func sanitizePeriod(period string) (Period, error) {
	var cleanPeriod string
	switch {
	case strings.Contains(strings.ToLower(period), "today"):
		cleanPeriod = time.Now().UTC().Weekday().String()
	case strings.Contains(strings.ToLower(period), "tonight"):
		cleanPeriod = fmt.Sprintf("%s night", time.Now().UTC().Weekday().String())
	case strings.Contains(strings.ToLower(period), "tomorrow night"):
		cleanPeriod = fmt.Sprintf("%s night", time.Now().UTC().AddDate(0, 0, 1).Weekday().String())
	case strings.Contains(strings.ToLower(period), "tomorrow"):
		cleanPeriod = time.Now().UTC().AddDate(0, 0, 1).Weekday().String()
	default:
		cleanPeriod = strings.ToLower(period)
	}
	m := periodRE.FindStringSubmatch(separatorRE.ReplaceAllString(cleanPeriod, " "))
	if m != nil && m[2] == "" {
		return Period{
			asKey:     strings.ToLower(m[1]),
			asName:    strings.Title(m[1]),
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
		return Period{}, errors.New("invalid period")
	}
}

// RandomPeriod generates a random day of the week and time of day
// and returns the corresponding Period struct
func RandomPeriod() Period {
	dayOfWeek := time.Now().UTC().AddDate(0, 0, rand.Intn(7)).Weekday().String()
	timeOfDay := timesOfDay[rand.Intn(2)]
	p, _ := sanitizePeriod(fmt.Sprintf("%s%s", dayOfWeek, timeOfDay))
	return p
}
