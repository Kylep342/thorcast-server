#!/usr/bin/env python
import argparse
import os

import geocode
import forecast as fc


def thorcast(city, state, api_key):
    coordinates = geocode.geocode(city, state, api_key)
    forecast = fc.get_forecast(**coordinates)
    return forecast