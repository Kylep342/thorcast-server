FROM python:3.7.2

# ubuntu setup
# RUN apt-get update -y
# RUN apt-get upgrade -y
# RUN apt-get install python -y 

# setup of thorcast
COPY . /app
WORKDIR /app
RUN pip install -r requirements.txt

# run thorcast
CMD python app.py