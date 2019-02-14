import argparse

import thorcast.geocode as geocode
import thorcast.forecast as fc
import utils.formatters as fmts


def lookup(city, state, period, thorcast_conn, redis_conn):
    """
    Main API function. Facilitates forecasting from request.

    Arguments:
        city:           [string]:
        state:          [string]:   
        thorcast_conn:  [sqlalchemy.engine.base.Connection]: DB conn
    """
    period = fmts.sanitize_period(period)
    city, state = fmts.sanitize_location(city, state)
    key = f'{city}_{state}_{period}'.lower().replace(' ', '_')
    forecast = redis_conn.lookup(key)
    if not forecast:
        coordinates = thorcast_conn.locate(city, state)
        if not coordinates:
            coordinates = geocode.geocode(city, state)
            thorcast_conn.register(city, state, **coordinates)
        forecasts_json = fc.forecast_from_api(**coordinates)
        forecasts = forecasts_json['properties']['periods']
        redis_conn.cache_forecasts(city, state, forecasts)
        forecast = redis_conn.lookup(key)
    return forecast


def deliver(forecast_json):
    return forecast_json['detailedForecast']