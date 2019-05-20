package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func (a *App) CacheForecasts(c City, s State, p Period, f Forecasts) string {
	now := time.Now()
	var detailedForecast string
	for _, forecast := range f.Properties.Periods {
		dayOfWeek := forecast.StartTime.Weekday().String()
		var timeOfDay string
		if forecast.IsDaytime {
			timeOfDay = ""
		} else {
			timeOfDay = "_night"
		}
		key := fmt.Sprintf(
			"%s_%s_%s%s",
			c.asKey,
			s.asKey,
			strings.ToLower(dayOfWeek),
			timeOfDay)
		err := a.Redis.Set(
			key,
			forecast.DetailedForecast,
			forecast.EndTime.Sub(now)).Err()
		if err != nil {
			log.Fatal(err)
		}
		if dayOfWeek == p.dayOfWeek && forecast.IsDaytime == p.isDaytime {
			detailedForecast = forecast.DetailedForecast
		}
	}
	return detailedForecast
}

func (a *App) LookupForecast(c City, s State, p Period) (string, error) {
	key := fmt.Sprintf(
		"%s_%s_%s",
		c.asKey,
		s.asKey,
		p.asKey)
	val, err := a.Redis.Get(key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
	
}