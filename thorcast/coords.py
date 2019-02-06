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
        insert_template = sqlalchemy.sql.text("""
            INSERT INTO geocodex (city, state, lat, lng, requests)
            VALUES (:city, :state, :lat, :lng, :requests)
        ;""")
        bind_params = {'city': city, 'state': state, 'lat': lat, 'lng': lng, 'requests': 1}
        self.conn.execute(insert_template, **bind_params) 
    
    def locate(self, city, state):
        query_template = sqlalchemy.sql.text("""
            SELECT lat, lng
            FROM geocodex
            WHERE city = :city
            AND state = :state
        ;""")
        bind_params = {'city': city, 'state': state}
        results = self.conn.execute(query_template, **bind_params)
        row = results.fetchone()
        if row:
            coordinates = {'lat': row['lat'], 'lng': row['lng']}
            update_template = sqlalchemy.sql.text("""
                UPDATE geocodex
                SET requests = requests + 1
                WHERE city = :city
                AND state = :state
            ;""")
            self.conn.execute(update_template, **bind_params)
            return coordinates
        else:
            return False