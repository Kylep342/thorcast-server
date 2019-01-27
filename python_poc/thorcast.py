import argparse
import os

import yaml

import geocode
import forecast as fc


with open(os.path.join(os.getcwd(), 'config.yml'), 'r') as conffile:
    config = yaml.safe_load(conffile)

def thorcast(city, state):
    coordinates = geocode.geocode(city, state, config['GoogleMapsAPIKey'])
    forecast = fc.get_forecast(**coordinates)
    return forecast


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('-c', '--city', help='Name of city you want forcast for', type=str)
    parser.add_argument('-s', '--state', help='Y coordinate of your current location', type=str)
    args = parser.parse_args()

    print(thorcast(args.city, args.state))

