import thorcast


def main():
    print('Welcome to Thorcast!')
    while True:
        address = input('Please enter the city and state, separated by a comma, to retrieve a forecast.\n')
        try:
            city, state = list(map(lambda x: x.strip(), address.split(',')))
        except ValueError:
            print('Please enter city, then state, separated by a comma (",")')
            continue
        print(thorcast.thorcast(city, state))
        another = input('Would you like to check another forecast? [y/n]\n')
        if another.lower() in ['n', 'no']:
            exit()
        else:
            continue


if __name__ == '__main__':
    main()