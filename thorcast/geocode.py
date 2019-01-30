import requests


def geocode(city, state, apikey):
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