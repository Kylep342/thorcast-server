"""
Main server script

Has Database connection setup, config loading, and routes

Author: Kyle Pekosh
Copyright 2019 by Kyle Pekosh
"""
import json
import logging
import os
import sys

from flask import Flask

import thorcast.thorcast as thorcast
import thorcast.geocodex as gx
import thorcast.weather_cache as wc


GC_USERNAME = os.getenv('THORCAST_DB_USERNAME')
GC_PASSWORD = os.getenv('THORCAST_DB_PASSWORD')
GC_HOST = os.getenv('THORCAST_DB_HOST')
GC_PORT = os.getenv('THORCAST_DB_PORT')
GC_DB = os.getenv('THORCAST_DB_NAME')

geocodex = gx.Geocodex(
    GC_USERNAME,
    GC_PASSWORD,
    GC_HOST,
    GC_PORT,
    GC_DB
)

REDIS_HOST = os.getenv('REDIS_HOST')
REDIS_PORT = os.getenv('REDIS_PORT')
REDIS_DB = os.getenv('REDIS_DB')
REDIS_PASSWORD = os.getenv('REDIS_PASSWORD')

weather_cache = wc.WeatherCache(
    REDIS_HOST,
    REDIS_PORT,
    REDIS_DB,
    REDIS_PASSWORD
)

root = logging.getLogger(__name__)

handler = logging.StreamHandler(sys.stdout)
formatter = logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s')
handler.setFormatter(formatter)
root.addHandler(handler)
root.setLevel(logging.INFO)




app = Flask(__name__)


@app.route('/')
def home():
    return('<html><body><h1>Welcome to Thorcast!</h1></body></html>')


@app.route('/api/city=<city>&state=<state>', defaults={'period': 'today'})
@app.route('/api/city=<city>&state=<state>&period=<period>')
def lookup(city, state, period):
    forecast_json = thorcast.lookup(
        city,
        state,
        period,
        geocodex,
        weather_cache,
        root
    )
    data = thorcast.deliver(city, state, period, forecast_json, root)
    response = app.response_class(
        response=json.dumps(data),
        status=200,
        mimetype='application/json'
    )
    return response


if __name__ == '__main__':
    app.run(host="0.0.0.0", port=5000, debug=True)
