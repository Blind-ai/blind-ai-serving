FROM python:3

ENV GPU=false
COPY . /app
WORKDIR /app
RUN pip install -r requirements.txt

CMD [ "python", "./serv-ai.py" ]