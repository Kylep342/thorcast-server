"""
This file contains functionality to interface with a database

The database is used to store coordinates of previously requested
locations. This is implemented to prevent going over GCP's free
allotment of 40k free geocode requests/month.
"""

import datetime
import os

import sqlalchemy

def gen_conn_str(username, password, host, port, db):

    return f'postgresql://{username}:{password}@{host}:{port}/{db}'


class Geocodex(object):
    """
    Geocodex is a database connection/interface object

    Since the app/model is small, this class represents the only table
    in the model.
    
    Methods include:
      - locate()    [look up coordinates for a location]
      - register()  [store coordinates for a location]
    """
    def __init__(self, username, password, host, port, db):
        """
        Establishes engine and connection to a database
        
        Args:
            username [string]: Database username to login
            password [string]: Database password to login
            host     [string]: Host/url for the database server
            port     [int]:    Port number the database listens through
            db       [string]: Name of the database
        """
        self.conn_str = gen_conn_str(username, password, host, port, db)
        self.engine = sqlalchemy.create_engine(self.conn_str)
        self.conn = self.engine.connect()

    def _exists(self):
        pass

    def _set_up(self):
        pass

    def register(self, city, state, lat, lng):
        """
        Method to encode the coordinates of a location to the db
        
        Args:
            city    [string]:   Name of the city being forecasted
            state   [string]:   US state that contains the city
            lat     [float]:    Latitude coordinate of the city
            lng     [float]:    Longitude coordinate of the city
        """
        insert_template = sqlalchemy.sql.text("""
            INSERT INTO geocodex (city, state, lat, lng, requests)
            VALUES (:city, :state, :lat, :lng, :requests)
        ;""")
        bind_params = {'city': city, 'state': state, 'lat': lat, 'lng': lng, 'requests': 1}
        self.conn.execute(insert_template, **bind_params) 
    
    def locate(self, city, state):
        """
        Method to retrieve the coordinates of a location in the db

        Args:
            city    [string]:   Name of the city being forecasted
            state   [string]:   US state that contains the city
        
        Returns:
            coordinates [dict{string:float}]:
                coordinate pair in the form: {'lat': , 'lng': }
            False       [boolean]: If the city has no record
        """
        query_template = sqlalchemy.sql.text("""
            SELECT lat, lng
            FROM geocodex
            WHERE LOWER(city) = LOWER(:city)
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