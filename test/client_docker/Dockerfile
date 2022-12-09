# syntax=docker/dockerfile:1

FROM python:3.8-alpine

WORKDIR /client

COPY requirements.txt requirements.txt
RUN pip3 install -r requirements.txt

COPY . .

CMD ["python", "client.py"]