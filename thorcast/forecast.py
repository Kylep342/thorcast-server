import requests


def forecast_from_api(lat, lng):
    """
    Function to get a forecast from api.weather.gov

    Arguments:
        lat [float]:    Latitude of city & state forecasting
        lng [float]:    Longitude of city & state forecasting
    
    Returns:
        forecast    [dict]: Object containing api.weather.gov response
    """
    try:
        points_resp = requests.get(f'https://api.weather.gov/points/{lat},{lng}')
        points_resp.raise_for_status()
        office = points_resp.json()

        #city = office['properties']['relativeLocation']['properties']['city']
        #state = office['properties']['relativeLocation']['properties']['state']

        forecast_endpt = office['properties']['forecast']

        forecast_resp = requests.get(forecast_endpt)
        forecast_resp.raise_for_status()
        forecast = forecast_resp.json()
    except Exception as e:
        raise e
    else:
        return forecast

