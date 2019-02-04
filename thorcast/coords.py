import datetime
import os

import sqlalchemy

def gen_conn_str(username, password, host, port, db):

    return f'postgresql://{username}:{password}@{host}:{port}/{db}'


class Geocodex(object):
    def __init__(self, username, password, host, port, db):
        self.conn_str = gen_conn_str(username, password, host, port, db)
        self.engine = sqlalchemy.create_engine(self.conn_str)
        self.conn = self.engine.connect()

    def _exists(self):
        pass

    def _set_up(self):
        pass

    def register(self, city, state, lat, lng):
        self.conn.execute("""
            INSERT INTO geocodex (
                city,
                state,
                lat,
                lng, 
                created_at
            ) VALUES (?, ?, ?, ?)
            ;""",
            city,
            state,
            lat,
            lng
        )
    
    def locate(self, city, state):
        rows = self.conn.execute("""
            SELECT
                lat,
                lng
            FROM geocodex
            WHERE city = ?
            AND state = ?
            ;""",
            city,
            state
        )
        results = rows.fetchall()
        if results:
            coordinates = {'lat': results[0], 'lng': results[1]}
        else:
            coordinates = False
        return coordinates