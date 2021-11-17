package cache

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis"

	"github.com/kylep342/thorcast-server/pkg/apis"
	"github.com/kylep342/thorcast-server/pkg/utils"
)

// CacheDetailedForecasts stores the provided forecasts
// for the given City, State, and Period
// key format is city.asKey_state.asKey_period.asKey
func CacheDetailedForecasts(
	cache *redis.Client
	city utils.City,
	state utils.State,
	period utils.Period,
	forecasts apis.Forecasts
) string {
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
		err := cache.Set(
			key,
			forecast.DetailedForecast,
			fcEndTime.Sub(now)).Err()
		if err != nil {
			log.Printf("Error occurred when setting a detailedForecast in Redis\nError is: %s\n", err.Error())
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
func LookupDetailedForecast(
	cache *redis.Client
	city utils.City,
	state utils.State,
	period utils.Period
) (string, error) {
	key := fmt.Sprintf(
		"%s_%s_%s",
		city.asKey,
		state.asKey,
		period.asKey)
	val, err := cache.Get(key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

// CacheHourlyForecasts persists all hourly forecasts in Redis as a list
// with an expiry of one hour
func CacheHourlyForecasts(
	cache *redis.Client
	city utils.City,
	state utils.State,
	hours int64,
	forecasts apis.Forecasts
) []string {
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
			"%s Forecast: %s, Temperature: %d\u00B0 %s, Wind: %s %s",
			fcDate.Format(time.RFC3339),
			fc.ShortForecast,
			int(fc.Temperature),
			fc.TemperatureUnit,
			fc.WindSpeed,
			fc.WindDirection)
		hourlyForecasts = append(hourlyForecasts, forecast)
	}
	err := cache.RPush(key, hourlyForecasts).Err()
	if err != nil {
		log.Printf("Error occurred when appending an hourly forecast to a list\nError is: %s\n", err.Error())
	}
	err = cache.ExpireAt(key, expiry).Err()
	if err != nil {
		log.Printf("Error occurred when setting an expiry for a list\nError is: %s\n", err.Error())
	}
	return hourlyForecasts[:hours]
}

// LookupHourlyForecast checks Redis for the requested city, state pair over
// the given hour range
// If a key is found, it returns the requested number of hourly forecasts
func LookupHourlyForecast(
	cache *redis.Client,
	city utils.City,
	state utils.State,
	hours int64
) ([]string, error) {
	key := fmt.Sprintf(
		"%s_%s_hourly",
		city.asKey,
		state.asKey)
	val, err := cache.LRange(key, 0, hours-1).Result()
	if err != nil {
		log.Printf("Error occurred when reading hourly forecasts from a list\nError is: %s\n", err.Error())
		return []string{}, err
	}
	// len(val) == 0 means key does not exist
	if len(val) > 0 {
		return val, nil
	}
	return []string{}, redis.Nil
}
