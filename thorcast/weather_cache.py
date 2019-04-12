"""
Cache interface utility


"""
import datetime
import json

import redis

import utils.calendar as cal


class WeatherCache(object):
    def __init__(self, host, port, db, password):
        self.conn = redis.Redis(host=host, port=port, db=db, password=password)
    
    def cache(self, key, value):
        # Get time at beginning of transaction
        # Use time to compute expiry of key (beginning of the next day)
        ts = datetime.datetime.now()
        return self.conn.set(
            key,
            json.dumps(value),
            ex=(86400 - (ts.hour * 3600 + ts.minute * 60 + ts.second))
        )
    
    def lookup(self, key):
        try:
            value = json.loads(self.conn.get(key))
        except TypeError:
            value = False
        except Exception as e:
            raise e
        else:
            return value
    
    def cache_forecasts(self, city, state, forecasts):
        for forecast in forecasts:
            try:
                dt_str = forecast['startTime']
                dt = datetime.datetime.strptime(dt_str, '%Y-%m-%dT%H:%M:%S%z')
                dayname = cal.day_of_week(dt)
                suffix = '_night' if not forecast['isDaytime'] else ''
                key = f'{city}_{state}_{dayname}{suffix}'.lower().replace(' ', '_')
                self.cache(key, forecast)
            except Exception as e:
                raise e
