import argparse

import requests


def clean_forecast(forecast):
    return forecast.replace('. ', '.\n')


def get_forecast(x, y):
    try:
        points_resp = requests.get(f'https://api.weather.gov/points/{x},{y}')
        points_resp.raise_for_status()
        office = points_resp.json()

        city = office['properties']['relativeLocation']['properties']['city']
        state = office['properties']['relativeLocation']['properties']['state']

        forecast_endpt = office['properties']['forecast']

        forecast_resp = requests.get(forecast_endpt)
        forecast_resp.raise_for_status()
        forecast = forecast_resp.json()

        forecast_p0 = forecast['properties']['periods'][0]
    except Exception as e:
        raise e
    finally:
        return f'{forecast_p0["name"]}\'s forecast for {city}, {state}:\n{clean_forecast(forecast_p0["detailedForecast"])}'


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('-x', help='X coordinate of your current location', type=float)
    parser.add_argument('-y', help='Y coordinate of your current location', type=float)
    args = parser.parse_args()

    print(get_forecast(args.x, args.y))

