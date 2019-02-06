#!/usr/bin/env python
import argparse
import os

import thorcast.geocode as geocode
import thorcast.forecast as fc
import thorcast.coords as coords


USERNAME = os.getenv('THORCAST_DB_USERNAME')
PASSWORD = os.getenv('THORCAST_DB_PASSWORD')
HOST = os.getenv('THORCAST_DB_HOST')
PORT = os.getenv('THORCAST_DB_PORT')
DB = os.getenv('THORCAST_DB_NAME')


geocodex = coords.Geocodex(USERNAME, PASSWORD, HOST, PORT, DB)


def lookup(city, state):
    coordinates = geocodex.locate(city, state)
    if not coordinates:
        coordinates = geocode.geocode(city, state)
        geocodex.register(city, state, **coordinates)
    forecast = fc.get_forecast(**coordinates)
    return forecast