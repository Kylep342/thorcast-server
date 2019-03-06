# Thorcast-Server

[![Join the chat at https://gitter.im/thorcast/community](https://badges.gitter.im/thorcast/community.svg)](https://gitter.im/thorcast/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Thorcast is a Discord/Slack chatbot that provides weather forecasts on demand

## Getting started
- Clone the repo at https://github.com/Kylep342/thorcast-server.git
- Run tests: `cd tests; pytest`
- Set up your own 'env.list' file in the structure demonstrated in env.list.example
- Run `docker build -t kylep342/thorcast-server .` to set up the Docker Image

## Usage
`docker run --env-file env.list -p 5000:5000 kylep342/thorcast-server` will start the server at 0.0.0.0:5000

Example:

```Bash
curl http://0.0.0.0:5000/api/city=Chicago&state=IL

{"forecast": "Today's forecast for Chicago, IL\nMostly sunny, with a high near 28.\nWind chill values as low as -2.\nWest southwest wind 10 to 20 mph, with gusts as high as 30 mph."}
```

## Upcoming features
- Move to Rust/Go

