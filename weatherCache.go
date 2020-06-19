package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

// CacheDetailedForecasts stores the provided forecasts
// for the given City, State, and Period
// key format is city.asKey_state.asKey_period.asKey
func (a *App) CacheDetailedForecasts(city City, state State, period Period, forecasts Forecasts) string {
	now := time.Now().UTC()
	var detailedForecast string
	for _, forecast := range forecasts.Properties.Periods {
		// log.Printf("Start time is: %s, end time is: %s\n", forecast.StartTime, forecast.EndTime)
		fcStartTime, _ := time.Parse(time.RFC3339, forecast.StartTime)
		fcEndTime, _ := time.Parse(time.RFC3339, forecast.EndTime)
		// log.Printf("fcStartTime is: %v, fcEndTime is: %v\n", fcStartTime, fcEndTime)
		dayOfWeek := fcStartTime.Weekday().String()
		var timeOfDay string
		if forecast.IsDaytime {
			timeOfDay = ""
		} else {
			timeOfDay = "_night"
		}
		key := fmt.Sprintf(
			"%s_%s_%s%s",
			city.asKey,
			state.asKey,
			strings.ToLower(dayOfWeek),
			timeOfDay)
		// log.Printf("Key is %s\n", key)
		err := a.Redis.Set(
			key,
			forecast.DetailedForecast,
			fcEndTime.Sub(now)).Err()
		if err != nil {
			log.Fatal(err)
		}
		// log.Printf("dayOfWeek is %s, period DOW is %s. fcDayTime is %t, periodDayTime is %t\n", dayOfWeek, period.dayOfWeek, forecast.IsDaytime, period.isDaytime)
		if dayOfWeek == period.dayOfWeek && forecast.IsDaytime == period.isDaytime {
			detailedForecast = forecast.DetailedForecast
		}
	}
	return detailedForecast
}

// LookupDetailedForecast tries to retrieve the forecast from the cache
// for the given City, State, and Period
func (a *App) LookupDetailedForecast(city City, state State, period Period) (string, error) {
	key := fmt.Sprintf(
		"%s_%s_%s",
		city.asKey,
		state.asKey,
		period.asKey)
	val, err := a.Redis.Get(key).Result()
	if err != nil {
		return "", err
	}
	return val, nil

}

// CacheHourlyForecasts
func (a *App) CacheHourlyForecasts(city City, state State, hours int64, forecasts Forecasts) []string {
	key := fmt.Sprintf(
		"%s_%s_hourly",
		city.asKey,
		state.asKey)
	now := time.Now().UTC()
	expiry := now.Add(1 * time.Hour).Truncate(1 * time.Hour)
	var hourlyForecasts []string
	for _, fc := range forecasts.Properties.Periods {
		fcDate, _ := time.Parse(time.RFC3339, fc.StartTime)
		forecast := fmt.Sprintf(
			"%s Forecast: %s, %d\u00b0 %s, Wind: %s %s",
			fcDate.Format(time.RFC3339),
			fc.ShortForecast,
			fc.Temperature,
			fc.TemperatureUnit,
			fc.WindSpeed,
			fc.WindDirection)
		hourlyForecasts = append(hourlyForecasts, forecast)
	}
	err := a.Redis.RPush(key, hourlyForecasts).Err()
	if err != nil {
		log.Fatal(err)
	}
	err = a.Redis.ExpireAt(key, expiry).Err()
	if err != nil {
		log.Fatal(err)
	}
	return hourlyForecasts[:hours]
}

// LookupHourlyForecast
func (a *App) LookupHourlyForecast(city City, state State, hours int64) ([]string, error) {
	key := fmt.Sprintf(
		"%s_%s_hourly",
		city.asKey,
		state.asKey)
	val, err := a.Redis.LRange(key, 0, hours-1).Result()
	if err != nil {
		log.Fatal(err)
	}
	// len(val) == 0 means key does not exist
	if len(val) > 0 {
		return val, nil
	}
	return []string{}, redis.Nil
}

// exists, err := a.Redis.Exists(key).Result()
// if err != nil {
// 	log.Fatal(err)
// }
// if exists == 1 {
// 	val, err := a.Redis.LRange(key, 0, hours-1).Result()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return val, nil
// } else {
// 	return []string{}, redis.Nil
// }
// }
