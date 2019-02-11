"""
Main server script

Has Database connection setup, config loading, and routes

Author: Kyle Pekosh
Copyright 2019 by Kyle Pekosh
"""
import os

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

weather_cache = wc.WeatherCache(

)

app = Flask(__name__)


@app.route('/')
def home():
    return('<html><body><h1>Welcome to Thorcast!</h1></body></html>')


@app.route('/thorcast/city=<city>&state=<state>')
def lookup_forecast(city, state):
    return thorcast.lookup(city, state, geocodex)


if __name__ == '__main__':
    app.run(host="0.0.0.0", port=5000, debug=True)