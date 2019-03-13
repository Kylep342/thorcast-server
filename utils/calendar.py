import datetime


def day_of_week(dt):
    dow = dt.weekday()

    day_of_week_codex = {
        0: 'monday',
        1: 'tuesday',
        2: 'wednesday',
        3: 'thursday',
        4: 'friday',
        5: 'saturday',
        6: 'sunday'
    }

    return day_of_week_codex[dow]