#!/bin/bash

docker run --env-file env.list -p 127.0.0.1:5000:5000 kylep342/thorcast