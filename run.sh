#!/bin/bash

docker run --env-file env.list -it -p 127.0.0.1:5000:5000 thorcast