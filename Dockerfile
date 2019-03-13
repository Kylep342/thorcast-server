FROM python:3.7.2

# setup of thorcast
COPY . /app
WORKDIR /app
RUN pip install -r requirements.txt

RUN ["python", "-m", "pytest"]

EXPOSE 5000

# run thorcast
ENTRYPOINT ["python",  "server.py"]