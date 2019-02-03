from flask import Flask

import thorcast.thorcast as thorcast


app = Flask(__name__)


@app.route('/')
def home():
    return('<html><body><h1>Welcome to Thorcast!</h1></body></html>')


@app.route('/thorcast/city=<city>&state=<state>')
def lookup_forecast(city, state):
    return thorcast.lookup(city, state)


if __name__ == '__main__':
    app.run(host="0.0.0.0", port=5000, debug=True)