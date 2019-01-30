# Thorcast
Thorcast is currently a command line app, but will be a Discord/Slack chatbot that provides weather forecasts on demand

## Getting started
- Clone the repo at https://github.com/Kylep342/thorcast.git
- Set up your own 'config.yml' file in the structure demonstrated in config.yml.example
- Run `build.sh` to set up the Docker Image

## Usage
Example:
```bash
run.sh

Welcome to Thorcast!

Please enter the city and state, separated by a comma, to retrieve a forecast.
Chicago, IL



Tonight's forecast for Chicago, IL:
Patchy blowing snow and scattered snow showers.
Mostly cloudy.
Low around -21, with temperatures rising to around -19 overnight.
Wind chill values as low as -47.
West wind around 20 mph, with gusts as high as 35 mph.
Chance of precipitation is 30%.
```

## Upcoming features
- Migration from command line app to Discord
- Slack support
- Move to Rust/Go

