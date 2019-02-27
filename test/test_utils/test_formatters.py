import pytest

import utils.formatters as fmts

def test_sanitize_city_underscore():
    city = 'salt_lake_city'
    assert fmts._sanitize_city(city) == 'Salt Lake City'


def test_sanitize_city_html_escape():
    city = 'Fort+lauderdale'
    assert fmts._sanitize_city(city) == 'Fort Lauderdale'


def test_sanitize_city_single_word():
    city = 'Chicago'
    assert fmts._sanitize_city(city) == 'Chicago'


def test_sanitize_state_full_name():
    state = 'West Virginia'
    assert fmts._sanitize_state(state) == 'WV'


def test_sanitize_state_postal_code():
    state = 'mo'
    assert fmts._sanitize_state(state) == 'MO'


def test_sanitize_state_not_a_state():
    state = 'Spam'
    with pytest.raises(KeyError):
        fmts._sanitize_state(state)