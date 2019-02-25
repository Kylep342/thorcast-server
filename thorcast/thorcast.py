import time

import redis
import sqlalchemy

import thorcast.geocode as geocode
import thorcast.forecast as fc
import utils.formatters as fmts


def lookup(city, state, period, thorcast_conn, redis_conn, logger):
    """
    Main API function. Facilitates forecasting from request.

    Arguments:
        city:           [string]:       The city name to forcast
        state:          [string]:       The state hosting the city
        period:         [string]:       The day/time to forecast
        thorcast_conn:  [sqlalchemy.engine.base.Connection]: DB conn
        redis_conn:     [redis.Redis]:  Redis connection object
    """
    period = fmts.sanitize_period(period)
    city, state = fmts.sanitize_location(city, state)
    key = f'{city}_{state}_{period}'.lower().replace(' ', '_')

    redis_retries = 5
    while redis_retries:
        try:
            forecast = redis_conn.lookup(key)
            break
        except redis.exceptions.ConnectionError as e:
            logger.info('Disconnected from Redis. Attempting to reconnect...')
            logger.info(f'Attempt {6 - redis_retries}')
            if redis_retries == 0:
                logger.error('Connection to Redis lost')
                raise e
            redis_retries -= 1
            time.sleep(0.5)

    if not forecast:
        pg_retries = 5
        while pg_retries:
            try:
                coordinates = thorcast_conn.locate(city, state)
                break
            except sqlalchemy.exc.OperationalError as e:
                logger.info('Disconnected from Postgres. Attempting to reconnect...')
                logger.info(f'Attempt {6 - redis_retries}')
                if not pg_retries:
                    logger.error('Connection to Postgres lost')
                    raise e
                pg_retries -= 1
                time.sleep(0.5)
        if not coordinates:
            coordinates = geocode.geocode(city, state)
            thorcast_conn.register(city, state, **coordinates)
        forecasts_json = fc.forecast_from_api(**coordinates)
        forecasts = forecasts_json['properties']['periods']
        redis_conn.cache_forecasts(city, state, forecasts)
        forecast = redis_conn.lookup(key)
    else:
        thorcast_conn.increment(city, state)
    return forecast


def deliver(city, state, period, forecast_json, logger):
    period = period.replace('_', ' ').capitalize()
    forecast = forecast_json['detailedForecast'].replace('. ', '.\n')
    return {'forecast': f"{period}'s forecast for {city}, {state}" + '\n' + forecast}