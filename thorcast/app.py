import os

import thorcast


def load_config():
    return os.environ


def main():
    config = load_config()
    print('Welcome to Thorcast!')
    while True:
        address = input('\nPlease enter the city and state, separated by a comma, to retrieve a forecast.\n')
        try:
            city, state = list(map(lambda x: x.strip(), address.split(',')))
        except ValueError:
            print('Please enter city, then state, separated by a comma (",")\n')
            continue
        print('\n\n\n', thorcast.thorcast(city, state, config.get('GOOGLE_MAPS_API_KEY')), sep='', end='\n\n\n')
        another = input('\nWould you like to check another forecast? [y/n]\n')
        if another.lower() in ['n', 'no']:
            exit()
        else:
            continue


if __name__ == '__main__':
    main()