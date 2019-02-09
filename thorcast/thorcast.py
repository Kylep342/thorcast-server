import argparse

import thorcast.geocode as geocode
import thorcast.forecast as fc
import utils.formatters as fmts


def lookup(city, state, thorcast_conn):
    city, state = fmts.sanitize_location(city, state)
    coordinates = thorcast_conn.locate(city, state)
    if not coordinates:
        coordinates = geocode.geocode(city, state)
        thorcast_conn.register(city, state, **coordinates)
    forecast = fc.get_forecast(**coordinates)
    return forecast