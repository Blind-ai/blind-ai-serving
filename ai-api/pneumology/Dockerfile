FROM python:3

ENV GPU=false

RUN mkdir /app
COPY main.py /app
COPY requirements.txt /app
WORKDIR /app

RUN pip install -r requirements.txt

ENTRYPOINT [ "python", "-u", "./main.py" ]
