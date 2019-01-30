#!/usr/bin/env python
import argparse
import os

import thorcast.geocode as geocode
import thorcast.forecast as fc


def lookup(city, state, api_key):
    coordinates = geocode.geocode(city, state, api_key)
    forecast = fc.get_forecast(**coordinates)
    return forecast