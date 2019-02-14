"""
Cache interface utility


"""
import datetime
import json

import redis

import utils.calendar as clndr


class WeatherCache(object):
    def __init__(self, host, port, db, password):
        self.conn = redis.Redis(host=host, port=port, db=db, password=password)
    
    def cache(self, key, value):
        return self.conn.set(key, json.dumps(value))
    
    def lookup(self, key):
        try:
            value = json.loads(self.conn.get(key))
        except TypeError:
            value = False
        except Exception as e:
            raise e
        return value
    
    def cache_forecasts(self, city, state, forecasts):
        for forecast in forecasts:
            try:
                dt_str = forecast['startTime']
                dt = datetime.datetime.strptime(dt_str, '%Y-%m-%dT%H:%M:%S%z')
                dayname = clndr.day_of_week(dt)
                suffix = '_night' if not forecast['isDaytime'] else ''
                key = f'{city}_{state}_{dayname}{suffix}'.lower().replace(' ', '_')
                self.cache(key, forecast)
            except Exception as e:
                raise e
