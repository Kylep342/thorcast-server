"""
Main server script

Has Database connection setup, config loading, and routes

Author: Kyle Pekosh
Copyright 2019 by Kyle Pekosh
"""
import os

from flask import Flask

import thorcast.thorcast as thorcast
import thorcast.coords as coords


USERNAME = os.getenv('THORCAST_DB_USERNAME')
PASSWORD = os.getenv('THORCAST_DB_PASSWORD')
HOST = os.getenv('THORCAST_DB_HOST')
PORT = os.getenv('THORCAST_DB_PORT')
DB = os.getenv('THORCAST_DB_NAME')

geocodex = coords.Geocodex(USERNAME, PASSWORD, HOST, PORT, DB)

app = Flask(__name__)


@app.route('/')
def home():
    return('<html><body><h1>Welcome to Thorcast!</h1></body></html>')


@app.route('/thorcast/city=<city>&state=<state>')
def lookup_forecast(city, state):
    return thorcast.lookup(city, state, geocodex)


if __name__ == '__main__':
    app.run(host="0.0.0.0", port=5000, debug=True)