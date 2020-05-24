FROM tensorflow/tensorflow:1.15.2

ENV GPU=false

RUN mkdir /app
COPY main.py /app
COPY model_definitive.json /app
COPY model_definitive.h5 /app
COPY requirements.txt /app
WORKDIR /app

RUN python -m pip install --upgrade pip setuptools
RUN pip install -r requirements.txt

ENTRYPOINT [ "python", "./main.py" ]
