"""
"""
import re


def _sanitize_state(state):
    states = {
        'alabama': 'AL', 'al': 'AL',
        'alaska': 'AK', 'ak': 'AK',
        'arizona': 'AZ', 'az': 'AZ',
        'arkansas': 'AR', 'ar': 'AR',
        'california': 'CA', 'ca': 'CA',
        'colorado': 'CO', 'co': 'CO',
        'connecticut': 'CT', 'ct': 'CT',
        'delaware': 'DE', 'de': 'DE',
        'florida': 'FL', 'fl': 'FL',
        'georgia': 'GA', 'ga': 'GA',
        'hawaii': 'HI', 'hi': 'HI',
        'idaho': 'ID', 'id': 'ID',
        'illinois': 'IL', 'il': 'IL',
        'indiana': 'IN', 'in': 'IN',
        'iowa': 'IA', 'ia': 'IA',
        'kansas': 'KS', 'ks': 'KS',
        'kentucky': 'KY', 'ky': 'KY',
        'louisiana': 'LA', 'la': 'LA',
        'maine': 'ME', 'me': 'ME',
        'maryland': 'MD', 'md': 'MD',
        'massachusetts': 'MA', 'ma': 'MA',
        'michigan': 'MI', 'mi': 'MI',
        'minnesota': 'MN', 'mn': 'MN',
        'mississippi': 'MS', 'ms': 'MS',
        'missouri': 'MO', 'mo': 'MO',
        'montana': 'MT', 'mt': 'MT',
        'nebraska': 'NE', 'ne': 'NE',
        'nevada': 'NV', 'nv': 'NV',
        'new hampshire': 'NH', 'nh': 'NH',
        'new jersey': 'NJ', 'nj': 'NJ',
        'new mexico': 'NM', 'nm': 'NM',
        'new york': 'NY', 'ny': 'NY',
        'north carolina': 'NC', 'nc': 'NC',
        'north dakota': 'ND', 'nd': 'ND',
        'ohio': 'OH', 'oh': 'OH',
        'oklahoma': 'OK', 'ok': 'OK',
        'oregon': 'OR', 'or': 'OR',
        'pennsylvania': 'PA', 'pa': 'PA',
        'rhode island': 'RI', 'ri': 'RI',
        'south carolina': 'SC', 'sc': 'SC',
        'south dakota': 'SD', 'sd': 'SD',
        'tennessee': 'TN', 'tn': 'TN',
        'texas': 'TX', 'tx': 'TX',
        'utah': 'UT', 'ut': 'UT',
        'vermont': 'VT', 'vt': 'VT',
        'virginia': 'VA', 'va': 'VA',
        'washington': 'WA', 'wa': 'WA',
        'west virginia': 'WV', 'wv': 'WV',
        'wisconsin': 'WI', 'wi': 'WI',
        'wyoming': 'WY', 'wy': 'WY'
    }
    try:
        return states[state.lower()]
    except KeyError as e:
        raise e


def _sanitize_city(city):
    return re.sub("[^a-zA-Z ']+", ' ', city)


def sanitize_location(city, state):
    return _sanitize_city(city), _sanitize_state(state)