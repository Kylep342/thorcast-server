# Thorcast
Thorcast is currently a command line app, but will be a Discord/Slack chatbot that provides weather forecasts on demand

## Getting started
- Clone the repo at https://github.com/Kylep342/thorcast.git
- Set up a python virtual environment. **\*\*Please do this in a directory separate from this project\*\***
- Activate the virtual environment and install dependencies with 'pip install -r requirements.txt'
- Navigate to the 'python_poc' folder from the project root
- Set up your own 'config.yml' file in the structure demonstrated in config.yml.example

## Usage
Example command:
```bash
thorcast.py -c Chicago -s IL
```
Example output:
```bash
This Afternoon\'s forecast for Chicago, IL:
Partly sunny.
High near 12, with temperatures falling to around 9 in the afternoon.
South southwest wind around 5 mph.
```

## Upcoming features
- Migration from command line app to Discord
- Slack support
- Move to Rust

