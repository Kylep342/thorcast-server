import os

import requests


def geocode(city, state):
    """
    Function to retrieve coordinates for a city and state

    Arguments:
        city    [string]:   Name of the city to forecast
        state   [string]:   Name of the state hosting the city

    Returns:
        coorinates  [dict]: Coordinate pair of city and state
                            in the form {'lat':, 'lng':}
    """
    apikey = os.getenv('GOOGLE_MAPS_API_KEY')
    fmt_city = city.replace(' ', '+')
    fmt_state = state.upper()
    geocode_api = 'https://maps.googleapis.com/maps/api/geocode/json'
    try:
        geocode_resp = requests.get(f'{geocode_api}?address={fmt_city},{fmt_state}&key={apikey}')
        geocode_resp.raise_for_status()
        geocode = geocode_resp.json()
        coordinates = geocode['results'][0]['geometry']['location']
    except Exception as e:
        raise e
    else:
        return coordinates