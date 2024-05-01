FROM python:3.8

COPY ./main.py /usr/src/main.py
COPY ./router /usr/src/router
COPY ./config /usr/src/config
COPY ./celery_core /usr/src/celery_core
COPY ./service /usr/src/service

COPY ./requirements-app.txt /usr/src/

# RUN pip3 install --upgrade pip

RUN pip3 install -r /usr/src/requirements-app.txt

WORKDIR /usr/src

CMD gunicorn --bind 0.0.0.0:5000 main:app -w 4 -k uvicorn.workers.UvicornWorker --access-logfile - --error-logfile - --log-level info --timeout 1000
# CMD uvicorn main:app --port 5000
