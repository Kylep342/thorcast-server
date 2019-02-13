import datetime


def day_of_week(dt_str):
    dt = datetime.datetime.strptime(dt_str, '%Y-%m-%dT%H:%M:%S%z')
    dow = dt.weekday()

    dow_register = {
        0: 'monday',
        1: 'tuesday',
        2: 'wednesday',
        3: 'thursday',
        4: 'friday',
        5: 'saturday',
        6: 'sunday'
    }

    return dow_register[dow]