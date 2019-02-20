import datetime

import utils.calendar as clndr


def test_day_of_week():
    test_date = datetime.datetime(2019, 1, 1)
    assert clndr.day_of_week(test_date) == 'tuesday'
