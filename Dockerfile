FROM python:3.7.2

# setup of thorcast
ADD . /app
WORKDIR /app
RUN pip install -r requirements.txt

# run thorcast
CMD python app.py