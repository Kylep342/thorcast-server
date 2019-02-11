"""
Cache interface utility


"""
import datetime
import json

import redis


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
        today = datetime.datetime.utcnow()

