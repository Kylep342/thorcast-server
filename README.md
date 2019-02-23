# Thorcast-Server

[![Join the chat at https://gitter.im/thorcast/community](https://badges.gitter.im/thorcast/community.svg)](https://gitter.im/thorcast/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Thorcast is a Discord/Slack chatbot that provides weather forecasts on demand

## Getting started
- Clone the repo at https://github.com/Kylep342/thorcast-server.git
- Run tests: `cd tests; pytest`
- Set up your own 'env.list' file in the structure demonstrated in env.list.example
- Run `build.sh` to set up the Docker Image

## Usage
`run.sh` should start the server at 0.0.0.0:5000

Example:

```Bash
curl http://0.0.0.0:5000/thorcast/city=Chicago&state=IL


Tonight's forecast for Chicago, IL:
Patchy blowing snow and scattered snow showers.
Mostly cloudy.
Low around -21, with temperatures rising to around -19 overnight.
Wind chill values as low as -47.
West wind around 20 mph, with gusts as high as 35 mph.
Chance of precipitation is 30%.
```

## Upcoming features
- Slack support
- Move to Rust/Go

