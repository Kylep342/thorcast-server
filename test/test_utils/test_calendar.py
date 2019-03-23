import datetime

import utils.calendar as cal


def test_day_of_week():
    test_date = datetime.date(2019, 1, 1)
    assert cal.day_of_week(test_date) == 'tuesday'
