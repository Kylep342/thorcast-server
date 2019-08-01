# Thorcast-Server

[![Join the chat at https://gitter.im/thorcast/community](https://badges.gitter.im/thorcast/community.svg)](https://gitter.im/thorcast/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Thorcast is a Discord/Slack chatbot that provides weather forecasts on demand

## Getting started

- Clone the repo at https://github.com/Kylep342/thorcast-server.git
- Set up your own 'env.list' file in the structure demonstrated in env.list.example
- Run `docker build -t kylep342/thorcast-server .` to set up the Docker Image

## Usage

`docker run --env-file env.list -p 8000:8000 kylep342/thorcast-server` will start the server at 0.0.0.0:8000

Example:

```Bash
curl http://0.0.0.0:8000/api/forecast/city=Chicago&state=IL&period=Tuesday+night

{"city":"Chicago","forecast":"Showers and thunderstorms likely. Mostly cloudy, with a low around 59.","period":"Tuesday night","state":"IL"}
```

## Upcoming features

- Add tests in Go
