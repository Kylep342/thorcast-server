#!/usr/bin/env python
import argparse

import thorcast.geocode as geocode
import thorcast.forecast as fc


def lookup(city, state):
    coordinates = geocode.geocode(city, state)
    forecast = fc.get_forecast(**coordinates)
    return forecast